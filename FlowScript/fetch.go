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

type RewardData struct {
	Type         string
	NodeID       string
	DelegatorID  int
	Amount       float64
	EpochCounter uint64
	Timestamp    time.Time
	BlockHeight  uint64
}

type FlowClient struct {
	client *client.Client
	ctx    context.Context
}

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
	bFound := false
	blockStart := blockEnd - 249

	fmt.Println("Searching for rewards...")
	for !bFound {
		fmt.Print(".")

		// 搜尋兩種事件
		foundRewards, foundDelegator, err := fc.searchBothRewards(blockStart, blockEnd)
		if err != nil {
			return nil, err
		}

		allRewards = append(allRewards, foundRewards...)
		allRewards = append(allRewards, foundDelegator...)

		bFound = len(foundRewards) > 0 || len(foundDelegator) > 0

		if !bFound {
			blockStart -= 250
			blockEnd -= 250
		}
	}

	fmt.Printf("\nFound %d reward records\n", len(allRewards))
	return fc.filterRewards(allRewards), nil
}

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

		// 計算到星期三的天數
		daysToTargetDay := int(time.Wednesday - dayOfWeek)
		hoursToTargetHour := int(20 - dayOfHour)

		fmt.Printf("Block probe = %d, Timestamp.Weekday = %s, daystoTargetDay = %d, Hour of day = %d, Timestamp = %s\n",
			blockEndProbe, dayOfWeek, daysToTargetDay, dayOfHour, timestamp.Format("2006-01-02 15:04:05"))

		if daysToTargetDay == 0 {
			if hoursToTargetHour < -1 {
				// 時間晚於 20:00，減少步長
				blockEndProbe -= iterationStep / 2
			} else if hoursToTargetHour == 0 || hoursToTargetHour == -1 {
				// 時間在 20:00 或 21:00，開始搜尋
				return blockEndProbe, nil
			} else {
				// 時間早於 20:00，繼續往前找
				blockEndProbe -= iterationStep
			}
		} else {
			// 不是星期三，繼續往前找
			blockEndProbe -= iterationStep
		}
	}
	return blockEndProbe, nil
}

func (fc *FlowClient) searchBothRewards(blockStart, blockEnd uint64) ([]RewardData, []RewardData, error) {
	// 搜尋 RewardsPaid 事件
	rewardsPaid, err := fc.searchRewards(blockStart, blockEnd, REWARD_PAID_EVENT)
	if err != nil {
		return nil, nil, err
	}

	// 搜尋 DelegatorRewardsPaid 事件
	delegatorRewards, err := fc.searchRewards(blockStart, blockEnd, DELEGATOR_REWARD_PAID_EVENT)
	if err != nil {
		return nil, nil, err
	}

	return rewardsPaid, delegatorRewards, nil
}

func (fc *FlowClient) searchRewards(blockStart, blockEnd uint64, eventType string) ([]RewardData, error) {
	var rewards []RewardData

	blockEvents, err := fc.client.GetEventsForHeightRange(fc.ctx, eventType, blockStart, blockEnd)
	if err != nil {
		return nil, err
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
				// 計算 epoch counter (簡化版本)
				epochCounter := fc.calculateEpochCounter(blockByID.Height)

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

				fmt.Printf("\nFound our node...block ID = %s, blockEnd = %d, amount = %.8f\n",
					blockEvent.BlockID.String(), blockEnd, amount)
				fmt.Printf("Block: %s Timestamp: %s\n",
					blockEvent.BlockID.String(), blockByID.Timestamp.String())
				if delegatorID > 0 {
					fmt.Printf("\tDelegatorID = %d\n", delegatorID)
				}
				fmt.Printf("\tEvent type: %s\n", event.Type)
				fmt.Printf("\tEvent: %s\n", event.Value.String())
			}
		}
	}

	return rewards, nil
}

func (fc *FlowClient) extractEventData(eventString string) (string, int, float64, error) {
	pattern := `nodeID: "([^"]+)", (?:\s*delegatorID: (\d+),)?\s*amount: ([0-9.]+)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(eventString)

	lenmatches := len(matches)

	if lenmatches < 3 {
		return "", -1, 0, fmt.Errorf("unable to extract node ID and amount from event string")
	}

	nodeID := strings.TrimSpace(matches[1])
	var delegatorID int
	if lenmatches == 4 && matches[2] != "" {
		delegatorID, _ = strconv.Atoi(strings.TrimSpace(matches[2]))
	} else {
		delegatorID = -1
	}

	amount, err := strconv.ParseFloat(matches[lenmatches-1], 64)
	if err != nil {
		return "", -1, 0, fmt.Errorf("unable to parse amount from event string: %w", err)
	}

	return nodeID, delegatorID, amount, nil
}

func (fc *FlowClient) getEventType(eventType string) string {
	if strings.Contains(eventType, "RewardsPaid") && !strings.Contains(eventType, "Delegator") {
		return "RewardsPaid"
	} else if strings.Contains(eventType, "DelegatorRewardsPaid") {
		return "DelegatorRewardsPaid"
	}
	return "Unknown"
}

func (fc *FlowClient) calculateEpochCounter(blockHeight uint64) uint64 {
	// 這是簡化的計算，實際情況可能需要更複雜的邏輯
	// 根據 Flow 的 epoch 機制，每個 epoch 大約包含特定數量的區塊
	// 這裡使用簡化的計算方式
	return blockHeight / 100000
}

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

func main() {
	// 印出開始時間
	currentTime := time.Now()
	timeStamp := currentTime.Format("2006-01-02 15:04:05")
	fmt.Println("Current time:", timeStamp)

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

	// 印出結束時間
	finishedTime := time.Now()
	timeStamp = finishedTime.Format("2006-01-02 15:04:05")
	fmt.Println("Finished time:", timeStamp)
	fmt.Println("***")
}
