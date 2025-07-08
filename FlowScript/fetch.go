package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	client "github.com/onflow/flow-go-sdk/access/grpc"
	grpc2 "google.golang.org/grpc"
)

const (
	TWM_NODE_ID                 = "e8f4bd649d08ecb5afb7023a0c5e8bb10ce56659399665da8cc9d517e7982f92"
	REWARD_PAID_EVENT           = "A.8624b52f9ddcd04a.FlowIDTableStaking.RewardsPaid"
	DELEGATOR_REWARD_PAID_EVENT = "A.8624b52f9ddcd04a.FlowIDTableStaking.DelegatorRewardsPaid"
	MAX_GRPC_MESSAGE_SIZE       = 1024 * 1024 * 16
	ACCESS_NODE_ADDRESS         = "access.mainnet.nodes.onflow.org:9000"
)

// RewardData 存儲獎勵資料的結構
type RewardData struct {
	Type         string
	NodeID       string
	DelegatorID  int
	Amount       float64
	EpochCounter uint64
	Timestamp    time.Time
	BlockHeight  uint64
}

// FlowClient 封裝 Flow 客戶端
type FlowClient struct {
	client *client.Client
	ctx    context.Context
}

// NewFlowClient 建立新的 Flow 客戶端
func NewFlowClient() (*FlowClient, error) {
	accessNodeClient, err := client.NewClient(
		ACCESS_NODE_ADDRESS,
		client.WithGRPCDialOptions(
			grpc2.WithInsecure(),
			grpc2.WithDefaultCallOptions(grpc2.MaxCallRecvMsgSize(MAX_GRPC_MESSAGE_SIZE)),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &FlowClient{
		client: accessNodeClient,
		ctx:    context.Background(),
	}, nil
}

// FetchThisWeekRewards 抓取本周的獎勵資料
func (fc *FlowClient) FetchThisWeekRewards() ([]RewardData, error) {
	var allRewards []RewardData

	// 取得最新區塊高度
	header, err := fc.client.GetLatestBlockHeader(fc.ctx, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block header: %w", err)
	}

	blockEnd := header.Height
	fmt.Printf("Latest block height: %d\n", blockEnd)

	// 跳轉到最近的星期三
	blockEnd, err = fc.jumpToWednesday(blockEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to jump to Wednesday: %w", err)
	}
	fmt.Printf("Wednesday block height: %d\n", blockEnd)

	// 搜尋最近的獎勵資料
	blockStart := blockEnd - 249
	found := false

	fmt.Println("Searching for rewards...")
	for !found {
		// 搜尋 RewardsPaid 事件
		rewardsPaid, foundRewards, err := fc.searchRewards(blockStart, blockEnd, REWARD_PAID_EVENT)
		if err != nil {
			return nil, err
		}
		allRewards = append(allRewards, rewardsPaid...)

		// 搜尋 DelegatorRewardsPaid 事件
		delegatorRewards, foundDelegator, err := fc.searchRewards(blockStart, blockEnd, DELEGATOR_REWARD_PAID_EVENT)
		if err != nil {
			return nil, err
		}
		allRewards = append(allRewards, delegatorRewards...)

		found = foundRewards || foundDelegator

		if !found {
			blockStart -= 250
			blockEnd -= 250
			fmt.Print(".")
		}
	}

	fmt.Printf("\nFound %d reward records\n", len(allRewards))
	return fc.filterRewards(allRewards), nil
}

// jumpToWednesday 跳轉到最近的星期三 20:00
func (fc *FlowClient) jumpToWednesday(blockEnd uint64) (uint64, error) {
	blockEndProbe := blockEnd
	iterationStep := uint64(5120)

	for blockEndProbe > 0 {
		blockByHeight, err := fc.client.GetBlockByHeight(fc.ctx, blockEndProbe)
		if err != nil {
			return blockEndProbe, err
		}

		timestamp := blockByHeight.Timestamp
		dayOfWeek := timestamp.Weekday()
		dayOfHour := timestamp.Hour()

		daysToTargetDay := int(time.Wednesday - dayOfWeek)
		hoursToTargetHour := int(20 - dayOfHour)

		fmt.Printf("Block: %d, Day: %s, Hour: %d, Timestamp: %s\n",
			blockEndProbe, dayOfWeek, dayOfHour, timestamp.Format("2006-01-02 15:04:05"))

		if daysToTargetDay == 0 {
			if hoursToTargetHour < -1 {
				blockEndProbe -= iterationStep / 2
			} else if hoursToTargetHour == 0 || hoursToTargetHour == -1 {
				return blockEndProbe, nil
			} else {
				blockEndProbe -= iterationStep
			}
		} else {
			blockEndProbe -= iterationStep
		}
	}
	return blockEndProbe, nil
}

// searchRewards 搜尋特定事件的獎勵資料
func (fc *FlowClient) searchRewards(blockStart, blockEnd uint64, eventType string) ([]RewardData, bool, error) {
	var rewards []RewardData
	found := false

	blockEvents, err := fc.client.GetEventsForHeightRange(fc.ctx, eventType, blockStart, blockEnd)
	if err != nil {
		return nil, false, err
	}

	for _, blockEvent := range blockEvents {
		// 取得區塊資訊以獲取時間戳
		blockByID, err := fc.client.GetBlockByID(fc.ctx, blockEvent.BlockID)
		if err != nil {
			continue
		}

		for _, event := range blockEvent.Events {
			nodeID, delegatorID, amount, err := fc.extractEventData(event.Value.String())
			if err != nil {
				continue
			}

			// 只處理我們的節點
			if nodeID == TWM_NODE_ID {
				found = true

				// 取得 epoch counter (這裡需要根據實際情況調整)
				epochCounter := fc.getEpochCounter(blockByID.Height)

				reward := RewardData{
					Type:         fc.getEventType(eventType),
					NodeID:       nodeID,
					DelegatorID:  delegatorID,
					Amount:       amount,
					EpochCounter: epochCounter,
					Timestamp:    blockByID.Timestamp,
					BlockHeight:  blockByID.Height,
				}
				rewards = append(rewards, reward)

				fmt.Printf("Found reward: Type=%s, DelegatorID=%d, Amount=%.8f\n",
					reward.Type, reward.DelegatorID, reward.Amount)
			}
		}
	}

	return rewards, found, nil
}

// extractEventData 從事件字串中提取資料
func (fc *FlowClient) extractEventData(eventString string) (string, int, float64, error) {
	pattern := `nodeID: "([^"]+)", (?:\s*delegatorID: (\d+),)?\s*amount: ([0-9.]+)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(eventString)

	if len(matches) < 3 {
		return "", -1, 0, fmt.Errorf("unable to extract data from event string")
	}

	nodeID := strings.TrimSpace(matches[1])
	var delegatorID int
	if len(matches) == 4 && matches[2] != "" {
		delegatorID, _ = strconv.Atoi(strings.TrimSpace(matches[2]))
	} else {
		delegatorID = -1
	}

	amount, err := strconv.ParseFloat(matches[len(matches)-1], 64)
	if err != nil {
		return "", -1, 0, fmt.Errorf("unable to parse amount: %w", err)
	}

	return nodeID, delegatorID, amount, nil
}

// getEventType 取得事件類型的簡化名稱
func (fc *FlowClient) getEventType(eventType string) string {
	if strings.Contains(eventType, "RewardsPaid") {
		return "RewardsPaid"
	} else if strings.Contains(eventType, "DelegatorRewardsPaid") {
		return "DelegatorRewardsPaid"
	}
	return "Unknown"
}

// getEpochCounter 取得 epoch counter (簡化實現)
func (fc *FlowClient) getEpochCounter(blockHeight uint64) uint64 {
	// 這裡需要根據 Flow 的實際邏輯來實現
	// 暫時返回一個計算值
	return blockHeight / 100000 // 簡化計算
}

// filterRewards 過濾獎勵資料，只保留需要的類型
func (fc *FlowClient) filterRewards(rewards []RewardData) []RewardData {
	var filtered []RewardData

	for _, reward := range rewards {
		// 只保留 RewardsPaid 或 DelegatorID 為 2、3、4 的 DelegatorRewardsPaid
		if reward.Type == "RewardsPaid" ||
			(reward.Type == "DelegatorRewardsPaid" &&
				(reward.DelegatorID == 2 || reward.DelegatorID == 3 || reward.DelegatorID == 4)) {
			filtered = append(filtered, reward)
		}
	}

	return filtered
}

// 主函數用於測試
func main() {
	client, err := NewFlowClient()
	if err != nil {
		panic(err)
	}

	rewards, err := client.FetchThisWeekRewards()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n=== Final Results ===\n")
	for _, reward := range rewards {
		fmt.Printf("Type: %s, NodeID: %s, DelegatorID: %d, Amount: %.8f, EpochCounter: %d, Timestamp: %s\n",
			reward.Type, reward.NodeID, reward.DelegatorID, reward.Amount, reward.EpochCounter, reward.Timestamp.Format("2006-01-02 15:04:05"))
	}
}
