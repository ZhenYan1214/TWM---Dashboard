/*
File: reward.go
Author: Jyh-Jong Liu
Created: 05/14/2023
Last Modified: 10/13/2023
Purpose: This tool retrieves block reward information from block events on the Flow blockchain network.
*/
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"encoding/base64"

	"github.com/onflow/flow-go-sdk"
	client "github.com/onflow/flow-go-sdk/access/grpc"
	grpc2 "google.golang.org/grpc" // required for increasing the max gRPC size
)

// twm node ID: e8f4bd649d08ecb5afb7023a0c5e8bb10ce56659399665da8cc9d517e7982f92
const TWM_NODE_ID = "e8f4bd649d08ecb5afb7023a0c5e8bb10ce56659399665da8cc9d517e7982f92"
const REWARD_PAID_EVENT = "A.8624b52f9ddcd04a.FlowIDTableStaking.RewardsPaid"
const DELEGATOR_REWARD_PAID_EVENT = "A.8624b52f9ddcd04a.FlowIDTableStaking.DelegatorRewardsPaid"

// increase max gRPC message size to 16 MB to resolve the issue for DelegatedRewardsPaid, default message size is 4 MB
const MAX_GRPC_MESSAGE_SIZE = 1024 * 1024 * 16

// An Identifier is a 32-byte unique identifier for an entity. From flow.go
type Identifier [32]byte

// RewardRecord struct 改為每一筆只帶自己的 delegator_total
// 其餘欄位維持原本格式
type RewardRecord struct {
	Type            string  `json:"type"`
	NodeID          string  `json:"node_id"`
	DelegatorID     int     `json:"delegator_id"`
	Amount          float64 `json:"amount"`
	EpochCounter    uint64  `json:"epoch_counter"`
	Timestamp       string  `json:"timestamp"`
	DelegatorTotal  string  `json:"delegator_total"`
	NodeTotal       string  `json:"node_total"`
	NodeStaked      string  `json:"node_staked"`
	DelegatorStaked string  `json:"delegator_staked"`
}

// block height is the height of the block containing the rewards payout events
func main() {
	// define command-line flags, JJ
	lastReward := flag.Bool("l", false, "Find last reward transactions.")
	allFlag := flag.Bool("a", false, "All reward transactions.")
	helpFlag := flag.Bool("h", false, "Help.")
	blockHeight := flag.Uint64("b", 0, "Block height ")
	totalB := flag.Uint64("t", 1, "Total blocks ")
	jumpStart := flag.Bool("j", false, "Jump start.")
	// Parse command-line flags
	flag.Parse()

	// fmt.Println("Arg = ", flag.NArg(), "Block height = ", *blockHeight)

	// Handle incorrect input and print usage
	if flag.NArg() != 0 || *helpFlag {
		fmt.Println("Usage: program_name [-h] [-l] [-a] [-j] [-b start_block_height] [-t total_blocks_to_process]")
		os.Exit(1)
	}

	var blockEnd uint64
	// 設定flow的節點
	accessNodeAddress := "access.mainnet.nodes.onflow.org:9000"

	// create a gRPC client for the Access node
	// accessNodeClient, err := client.NewClient(accessNodeAddress)
	// change the default gRPC message size
	accessNodeClient, err := client.NewClient(
		accessNodeAddress,
		client.WithGRPCDialOptions(
			grpc2.WithInsecure(),
			grpc2.WithDefaultCallOptions(grpc2.MaxCallRecvMsgSize(MAX_GRPC_MESSAGE_SIZE)),
		),
	)
	if err != nil {
		fmt.Println("err:", err.Error())
		panic(err)
	}

	ctx := context.Background()

	// print the time stamp
	currentTime := time.Now()
	// format the time stamp as "YYYY-MM-DD HH:MM:SS"
	timeStamp := currentTime.Format("2006-01-02 15:04:05")
	fmt.Println("Current time:", timeStamp)

	// 使用-b的話就會指定高度，沒有的話就是最新高度
	if *blockHeight != 0 {
		blockEnd = *blockHeight
	} else {
		// get latest block height
		/*
					   	    isSealed := true
					   	    latestBlock, err := accessNodeClient.GetLatestBlock(ctx, isSealed)
					   	    if err != nil {
			//		   		fmt.Println("err:", err.Error())
					   		panic(err)
					   	    }
					   // get the block by ID
					       	    blockID := latestBlock.ID.String()
					       	    blockByID, err := accessNodeClient.GetBlockByID(ctx, flow.HexToID(blockID))
					       	    printBlock(blockByID, err)
		*/
		header, err := accessNodeClient.GetLatestBlockHeader(ctx, true)
		if err != nil {
			fmt.Println("err:", err.Error())
			panic(err)
		}
		blockEnd = header.Height
	}

	fmt.Println("End height =", blockEnd, " All flag =", *allFlag)

	//看有沒有使用jumpstart
	// check if we want to jump start to Wed night 20:00
	if *jumpStart {
		fmt.Println("Jump to start height on Wednesday ...")
		blockEnd, err = JumpToNewEndHeight(ctx, accessNodeClient, blockEnd)
		if err != nil {
			panic(err)
		}
		fmt.Println("Find new end height =", blockEnd)
	}

	// 1. 先收集所有 delegator id
	type RawRecord struct {
		Type         string
		NodeID       string
		DelegatorID  int
		Amount       float64
		EpochCounter uint64
		Timestamp    string
	}
	var allRecordsRaw []RawRecord
	delegatorIDSet := make(map[int]struct{})
	blockStart := blockEnd - 249
	if *lastReward {
		bFound := false
		for bFound == false {
			records, err, found := GetBlockEventsCollectAll(ctx, blockStart, blockEnd, accessNodeClient, *allFlag, REWARD_PAID_EVENT, delegatorIDSet)
			if err != nil {
				panic(err)
			}
			for _, r := range records {
				allRecordsRaw = append(allRecordsRaw, r)
			}
			records, err, found2 := GetBlockEventsCollectAll(ctx, blockStart, blockEnd, accessNodeClient, *allFlag, DELEGATOR_REWARD_PAID_EVENT, delegatorIDSet)
			if err != nil {
				panic(err)
			}
			for _, r := range records {
				allRecordsRaw = append(allRecordsRaw, r)
			}
			bFound = found || found2
			blockStart -= 250
			blockEnd -= 250
		}
	} else {
		totalBlocks := *totalB
		iterations := int(totalBlocks / 250)
		var remaining uint64 = uint64(totalBlocks % 250)
		for i := 0; i < iterations; i++ {
			records, err, _ := GetBlockEventsCollectAll(ctx, blockStart, blockEnd, accessNodeClient, *allFlag, REWARD_PAID_EVENT, delegatorIDSet)
			if err != nil {
				panic(err)
			}
			for _, r := range records {
				allRecordsRaw = append(allRecordsRaw, r)
			}
			records, err, _ = GetBlockEventsCollectAll(ctx, blockStart, blockEnd, accessNodeClient, *allFlag, DELEGATOR_REWARD_PAID_EVENT, delegatorIDSet)
			if err != nil {
				panic(err)
			}
			for _, r := range records {
				allRecordsRaw = append(allRecordsRaw, r)
			}
			blockStart -= 250
			blockEnd -= 250
		}
		if remaining != 0 {
			blockStart += uint64(250) - remaining
			records, err, _ := GetBlockEventsCollectAll(ctx, blockStart, blockEnd, accessNodeClient, *allFlag, REWARD_PAID_EVENT, delegatorIDSet)
			if err != nil {
				panic(err)
			}
			for _, r := range records {
				allRecordsRaw = append(allRecordsRaw, r)
			}
			records, err, _ = GetBlockEventsCollectAll(ctx, blockStart, blockEnd, accessNodeClient, *allFlag, DELEGATOR_REWARD_PAID_EVENT, delegatorIDSet)
			if err != nil {
				panic(err)
			}
			for _, r := range records {
				allRecordsRaw = append(allRecordsRaw, r)
			}
		}
	}
	// 2. 查詢所有 delegator 的 tokensRewarded
	// 查詢 nodeTotal, nodeStaked
	nodeTotal, nodeStaked, err := GetNodeTotal()
	if err != nil {
		fmt.Println("查詢 node_total 失敗:", err)
		os.Exit(1)
	}
	// 查詢所有 delegator 的 tokensRewarded 和 tokensStaked
	delegatorTotals := queryDelegatorTotalsFixed(delegatorIDSet)
	// 3. 重新組裝 allRecords，帶上自己的 delegatorTotal
	var allRecords []RewardRecord
	for _, r := range allRecordsRaw {
		idStr := strconv.Itoa(r.DelegatorID)
		delegatorTotal := "0"
		delegatorStaked := "0"
		if v, ok := delegatorTotals[idStr]; ok {
			delegatorTotal = v.Rewarded
			delegatorStaked = v.Staked
		}
		rec := RewardRecord{
			Type:            r.Type,
			NodeID:          r.NodeID,
			DelegatorID:     r.DelegatorID,
			Amount:          r.Amount,
			EpochCounter:    r.EpochCounter,
			Timestamp:       r.Timestamp,
			DelegatorTotal:  delegatorTotal,
			NodeTotal:       nodeTotal,
			NodeStaked:      nodeStaked,
			DelegatorStaked: delegatorStaked,
		}
		allRecords = append(allRecords, rec)
	}
	// print struct
	fmt.Println("\n結果已經儲存成 struct，內容如下：")
	for _, r := range allRecords {
		fmt.Printf("%+v\n", r)
	}

	// 將結果生成 JSON 檔案
	if len(allRecords) > 0 {
		jsonData, err := json.MarshalIndent(allRecords, "", "  ") //轉成json，然後自動排版
		if err != nil {
			fmt.Printf("JSON marshal error: %v\n", err)
			os.Exit(1)
		}
		err = ioutil.WriteFile("rewards.json", jsonData, 0644)
		if err != nil {
			fmt.Printf("Write file error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ 已將 %d 筆資料存成 rewards.json\n", len(allRecords))
	}

	// print the finished time
	finishedTime := time.Now()
	// format the time stamp as "YYYY-MM-DD HH:MM:SS"
	timeStamp = finishedTime.Format("2006-01-02 15:04:05")
	fmt.Println("Finished time:", timeStamp)
	fmt.Println("***")
}

// 新增：收集所有 delegator id 的版本，並回傳原始欄位
func GetBlockEventsCollectAll(ctx context.Context, blockStart, blockEnd uint64, accessNodeClient *client.Client, all bool, eventString string, allDelegatorIDs map[int]struct{}) ([]struct {
	Type         string
	NodeID       string
	DelegatorID  int
	Amount       float64
	EpochCounter uint64
	Timestamp    string
}, error, bool) {
	fmt.Print(".")
	blockEvents, err := accessNodeClient.GetEventsForHeightRange(ctx, eventString, blockStart, blockEnd)
	bFound := false
	var results []struct {
		Type         string
		NodeID       string
		DelegatorID  int
		Amount       float64
		EpochCounter uint64
		Timestamp    string
	}
	if err != nil {
		panic(err)
	}
	for _, blockEvent := range blockEvents {
		for _, event := range blockEvent.Events {
			if all {
				fmt.Println("\tEvent type: " + event.Type)
				fmt.Println("\tEvent: " + event.Value.String())
				fmt.Println("\tEvent payload: " + string(event.Payload))
			} else {
				nodeID, delegatorID, amount, epochCounter, err := extractNodeIDAndAmount(event.Value.String())
				if err != nil {
					fmt.Println("Failed to extract node ID and amount:", err)
					return nil, err, bFound
				}
				if nodeID == TWM_NODE_ID {
					// 不再篩選 delegatorID，全部都收集
					blockByID, err := accessNodeClient.GetBlockByID(ctx, blockEvent.BlockID)
					if err != nil {
						fmt.Println("err:", err.Error())
						return nil, err, bFound
					}
					eventType := "Node"
					if strings.Contains(event.Type, "DelegatorRewardsPaid") {
						eventType = "Delegator"
					}
					// 收集 delegatorID
					allDelegatorIDs[delegatorID] = struct{}{}
					record := struct {
						Type         string
						NodeID       string
						DelegatorID  int
						Amount       float64
						EpochCounter uint64
						Timestamp    string
					}{
						Type:         eventType,
						NodeID:       nodeID,
						DelegatorID:  delegatorID,
						Amount:       amount,
						EpochCounter: epochCounter,
						Timestamp:    blockByID.Timestamp.Format("2006-01-02 15:04:05"),
					}
					results = append(results, record)
					bFound = true
				}
			}
		}
	}
	return results, err, bFound
}

// added time stamp on 7/19/2023
func PrintTimeStamp(ctx context.Context, accessNodeClient *client.Client, blockEvent flow.BlockEvents) error {
	// get the block by ID to retrieve the time stamp
	blockByID, err := accessNodeClient.GetBlockByID(ctx, blockEvent.BlockID)
	if err != nil {
		fmt.Println("err:", err.Error())
		return err
	}
	fmt.Println("Block: " + blockEvent.BlockID.String() + " Timestamp: " + blockByID.Timestamp.String())
	return nil
}

// Try to check if the timestamp information from block information is
// on Wednesday, otherwise loop back until find the block with timestamp
// that is on Wednesday
func JumpToNewEndHeight(ctx context.Context, accessNodeClient *client.Client, blockEnd uint64) (uint64, error) {
	blockEndProbe := blockEnd

	iterationStep := uint64(5120)

	for blockEndProbe > 0 {

		// The reward enents happens around 19:00 Wed UTC, we may improve the performance by
		// jumping the block with the time stamp 20:00 Wed UTC, if the time stamp
		/*
		   blockID := flow.NewHeight(blockEndProbe)

		   fmt.Printf("Block Height: %d\n", blockEndProbe)
		   fmt.Printf("Block ID: %s\n", blockID.String())

		   blockByID, err := accessNodeClient.GetBlockByID(ctx, blockID)
		*/
		blockByHeight, err := accessNodeClient.GetBlockByHeight(ctx, blockEndProbe)

		if err != nil {
			fmt.Println("err:", err.Error())
			return blockEndProbe, err
		}
		// the time stamp returned is UCT time zone
		// targeted time for reward paid is on Wed 19:00, we should skip blocks untile Wed 20:00 and then
		// check detail event information
		timestamp := blockByHeight.Timestamp

		dayOfWeek := timestamp.Weekday()
		dayOfHour := timestamp.Hour()

		// calculate the days to Wednesday (day value 3)
		daysToTargetDay := int(time.Wednesday - dayOfWeek)
		hoursToTargetHour := int(20 - dayOfHour)

		fmt.Println("Block probe = ", blockEndProbe, "Timestamp.Weekday = ", dayOfWeek, " daystoTargetDay = ", daysToTargetDay,
			"Hour of day = ", dayOfHour, "Timestamp =", timestamp)

		if daysToTargetDay == 0 {
			if hoursToTargetHour < -1 {
				// half of the steps if the time is later than 20:00
				blockEndProbe -= iterationStep / 2
			} else if hoursToTargetHour == 0 || hoursToTargetHour == -1 {
				// start from 20:00 or 21:00 and check events backward
				return blockEndProbe, nil
			} else {
				// in this case, the next reward is not paid yet, loop back to previous reward information
				blockEndProbe -= iterationStep
			}
		} else {
			// different day, need to skip
			blockEndProbe -= iterationStep
		}
	}
	return blockEndProbe, nil
}

// extract node ID, delegatorID, amount, and epochCounter from the event
func extractNodeIDAndAmount(eventString string) (string, int, float64, uint64, error) {
	// Regular expression pattern to match node ID, delegatorID, amount, and epochCounter
	// pattern := `nodeID: "([^"]+)", (?:\s*delegatorID: (\d+),)?\s*amount: ([0-9.]+)(?:,\s*epochCounter: (\d+))?`
	pattern := `nodeID: "([^"]+)", (?:\s*delegatorID: (\d+),)?\s*amount: ([0-9.]+)(?:,\s*epochCounter: (\d+))?`

	// Compile the regular expression pattern
	re := regexp.MustCompile(pattern)

	// Find the matches in the event string
	matches := re.FindStringSubmatch(eventString)

	lenmatches := len(matches)

	if lenmatches < 3 {
		return "", -1, 0, 0, fmt.Errorf("unable to extract node ID and amount from event string")
	}

	nodeID := strings.TrimSpace(matches[1])
	var delegatorID int
	if lenmatches >= 4 && matches[2] != "" {
		delegatorID, _ = strconv.Atoi(strings.TrimSpace(matches[2]))
	} else {
		delegatorID = -1
	}

	amount, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return "", -1, 0, 0, fmt.Errorf("unable to parse amount from event string: %w", err)
	}

	var epochCounter uint64
	if lenmatches >= 5 && matches[4] != "" {
		epochCounter, _ = strconv.ParseUint(strings.TrimSpace(matches[4]), 10, 64)
	} else {
		epochCounter = 0
	}

	return nodeID, delegatorID, amount, epochCounter, nil
}

// 新增 function 查詢 node_total (tokensRewarded)
func GetNodeTotal() (string, string, error) {
	url := "https://rest-mainnet.onflow.org/v1/scripts?blocks_height=final"
	method := "POST"
	payload := `{
		"script": "aW1wb3J0IEZsb3dJRFRhYmxlU3Rha2luZyBmcm9tIDB4ODYyNGI1MmY5ZGRjZDA0YQ0KDQovLyBUaGlzIHNjcmlwdCBnZXRzIGFsbCB0aGUgaW5mbyBhYm91dCBhIG5vZGUgYW5kIHJldHVybnMgaXQNCmFjY2VzcyhhbGwpIGZ1biBtYWluKG5vZGVJRDogU3RyaW5nKTogRmxvd0lEVGFibGVTdGFraW5nLk5vZGVJbmZvIHsNCiAgICByZXR1cm4gRmxvd0lEVGFibGVTdGFraW5nLk5vZGVJbmZvKG5vZGVJRDogbm9kZUlEKQ0KfQ==",
		"arguments": [
			"eyJ0eXBlIjoiU3RyaW5nIiwidmFsdWUiOiJlOGY0YmQ2NDlkMDhlY2I1YWZiNzAyM2EwYzVlOGJiMTBjZTU2NjU5Mzk5NjY1ZGE4Y2M5ZDUxN2U3OTgyZjkyIn0="
		]
	}`
	body1 := strings.NewReader(payload)
	req, err := http.NewRequest(method, url, body1)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}
	cleaned := strings.Trim(string(body), "\"")
	decoded, err := base64.StdEncoding.DecodeString(cleaned)
	if err != nil {
		return "", "", err
	}
	var outer struct {
		Value struct {
			Fields []struct {
				Name  string          `json:"name"`
				Value json.RawMessage `json:"value"`
			} `json:"fields"`
		} `json:"value"`
	}
	if err := json.Unmarshal(decoded, &outer); err != nil {
		return "", "", err
	}
	var nodeTotal, nodeStaked string
	for _, field := range outer.Value.Fields {
		if field.Name == "tokensRewarded" {
			var v struct {
				Value string `json:"value"`
			}
			if err := json.Unmarshal(field.Value, &v); err == nil {
				nodeTotal = v.Value
			}
		}
		if field.Name == "tokensStaked" {
			var v struct {
				Value string `json:"value"`
			}
			if err := json.Unmarshal(field.Value, &v); err == nil {
				nodeStaked = v.Value
			}
		}
	}
	return nodeTotal, nodeStaked, nil
}

// 查詢單一 delegator 的 tokensRewarded
// 查詢單一 delegator 的 tokensRewarded
func GetDelegatorTotal(delegatorID int) (string, string, error) {
	// 如果 delegatorID 是 -1，直接回傳 "0"，因為這代表 Node 本身而非 delegator
	if delegatorID == -1 {
		return "0", "0", nil
	}

	url := "https://rest-mainnet.onflow.org/v1/scripts?blocks_height=final"
	method := "POST"
	arg := fmt.Sprintf(`{"type":"UInt32","value":"%d"}`, delegatorID)
	argBase64 := base64.StdEncoding.EncodeToString([]byte(arg))
	payload := fmt.Sprintf(`{
        "script": "aW1wb3J0IEZsb3dJRFRhYmxlU3Rha2luZyBmcm9tIDB4ODYyNGI1MmY5ZGRjZDA0YQ0KDQovLyBUaGlzIHNjcmlwdCByZXR1cm5zIGFsbCB0aGUgaW5mbyBhc3NvY2lhdGVkIHdpdGggYSBkZWxlZ2F0b3INCg0KYWNjZXNzKGFsbCkgZnVuIG1haW4obm9kZUlEOiBTdHJpbmcsIGRlbGVnYXRvcklEOiBVSW50MzIpOiBGbG93SURUYWJsZVN0YWtpbmcuRGVsZWdhdG9ySW5mbyB7DQogICAgcmV0dXJuIEZsb3dJRFRhYmxlU3Rha2luZy5EZWxlZ2F0b3JJbmZvKG5vZGVJRDogbm9kZUlELCBkZWxlZ2F0b3JJRDogZGVsZWdhdG9ySUQpDQp9DQo=",
        "arguments": [
            "eyJ0eXBlIjoiU3RyaW5nIiwidmFsdWUiOiJlOGY0YmQ2NDlkMDhlY2I1YWZiNzAyM2EwYzVlOGJiMTBjZTU2NjU5Mzk5NjY1ZGE4Y2M5ZDUxN2U3OTgyZjkyIn0=",
            "%s"
        ]
    }`, argBase64)

	body1 := strings.NewReader(payload)
	req, err := http.NewRequest(method, url, body1)
	if err != nil {
		return "", "", fmt.Errorf("建立 HTTP 請求失敗: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("HTTP 請求失敗: %w", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", fmt.Errorf("讀取回應失敗: %w", err)
	}

	// 先檢查 HTTP 狀態碼
	if res.StatusCode != 200 {
		return "", "", fmt.Errorf("API 回傳錯誤狀態碼 %d: %s", res.StatusCode, string(body))
	}

	// 檢查回應是否為空
	if len(body) == 0 {
		return "", "", fmt.Errorf("API 回傳空回應")
	}

	// 移除前後的引號
	cleaned := strings.Trim(string(body), "\"")

	// 檢查是否為有效的 base64
	if len(cleaned) == 0 {
		return "", "", fmt.Errorf("清理後的回應為空")
	}

	// 嘗試 base64 解碼，並提供更詳細的錯誤訊息
	decoded, err := base64.StdEncoding.DecodeString(cleaned)
	if err != nil {
		return "", "", fmt.Errorf("base64 解碼失敗 (原始回應: %s): %w", string(body), err)
	}

	// 解析 JSON 結構
	var outer struct {
		Value struct {
			Fields []struct {
				Name  string `json:"name"`
				Value struct {
					Value string `json:"value"`
				} `json:"value"`
			} `json:"fields"`
		} `json:"value"`
	}

	if err := json.Unmarshal(decoded, &outer); err != nil {
		return "", "", fmt.Errorf("JSON 解析失敗 (解碼後內容: %s): %w", string(decoded), err)
	}

	// 尋找 tokensRewarded 欄位
	var tokensRewarded, tokensStaked string
	for _, field := range outer.Value.Fields {
		if field.Name == "tokensRewarded" {
			tokensRewarded = field.Value.Value
		}
		if field.Name == "tokensStaked" {
			tokensStaked = field.Value.Value
		}
	}
	return tokensRewarded, tokensStaked, nil
}

// 同時也修正主函數中的 delegator 總額查詢部分
func queryDelegatorTotalsFixed(delegatorIDSet map[int]struct{}) map[string]struct{ Rewarded, Staked string } {
	delegatorTotals := make(map[string]struct{ Rewarded, Staked string })

	for id := range delegatorIDSet {
		if id == -1 {
			delegatorTotals[strconv.Itoa(id)] = struct{ Rewarded, Staked string }{Rewarded: "0", Staked: "0"}
			continue
		}
		rewarded, staked, err := GetDelegatorTotal(id)
		if err != nil {
			fmt.Printf("查詢 delegator %d 總額失敗: %v\n", id, err)
			delegatorTotals[strconv.Itoa(id)] = struct{ Rewarded, Staked string }{Rewarded: "0", Staked: "0"}
		} else {
			delegatorTotals[strconv.Itoa(id)] = struct{ Rewarded, Staked string }{Rewarded: rewarded, Staked: staked}
		}
		// 添加小延遲，避免 API 請求過於頻繁
		time.Sleep(100 * time.Millisecond)
	}

	return delegatorTotals
}

/* sample reward paid message
	Event type: A.8624b52f9ddcd04a.FlowIDTableStaking.RewardsPaid
	Event: A.8624b52f9ddcd04a.FlowIDTableStaking.RewardsPaid(nodeID: "5a4bff17941a73909472afe23f1ccdc59d7526f93b16b4e374bd8353f8b624b4", amount: 9223.47334447)
	Event payload: {"value":{"id":"A.8624b52f9ddcd04a.FlowIDTableStaking.RewardsPaid","fields":[{"value":{"value":"5a4bff17941a73909472afe23f1ccdc59d7526f93b16b4e374bd8353f8b624b4","type":"String"},"name":"nodeID"},{"value":{"value":"9223.47334447","type":"UFix64"},"name":"amount"}]},"type":"Event"}

sample delegator reward paid message
	Event type: A.8624b52f9ddcd04a.FlowIDTableStaking.DelegatorRewardsPaid
	Event: A.8624b52f9ddcd04a.FlowIDTableStaking.DelegatorRewardsPaid(nodeID: "95ffacf0c05757cff71a4ee49e025d5a6d1103a3aa7d91253079e1bfb7c22458", delegatorID: 23, amount: 0.10424555)
	Event payload: {"value":{"id":"A.8624b52f9ddcd04a.FlowIDTableStaking.DelegatorRewardsPaid","fields":[{"value":{"value":"95ffacf0c05757cff71a4ee49e025d5a6d1103a3aa7d91253079e1bfb7c22458","type":"String"},"name":"nodeID"},{"value":{"value":"23","type":"UInt32"},"name":"delegatorID"},{"value":{"value":"0.10424555","type":"UFix64"},"name":"amount"}]},"type":"Event"}

*/
