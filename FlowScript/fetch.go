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

// 確保只保留這一份 RewardRecord struct 定義，其餘全部刪除
type RewardRecord struct {
	Type            string  `json:"type"`
	NodeID          string  `json:"node_id"`
	DelegatorID     int     `json:"delegator_id"`
	Amount          float64 `json:"amount"`
	EpochCounter    uint64  `json:"epoch_counter"`
	Timestamp       string  `json:"timestamp"`
	DelegatorTotal2 string  `json:"delegator_total2"`
	DelegatorTotal3 string  `json:"delegator_total3"`
	DelegatorTotal4 string  `json:"delegator_total4"`
	Node_total      string  `json:"node_total"`
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

	//看有沒有使用 -l ，一直找，直到找到reward事件
	var allRecords []RewardRecord
	blockStart := blockEnd - 249
	if *lastReward {
		bFound := false
		for bFound == false {
			records, err, found := GetBlockEvents(ctx, blockStart, blockEnd, accessNodeClient, *allFlag, REWARD_PAID_EVENT)
			if err != nil {
				panic(err)
			}
			allRecords = append(allRecords, records...)
			records, err, found2 := GetBlockEvents(ctx, blockStart, blockEnd, accessNodeClient, *allFlag, DELEGATOR_REWARD_PAID_EVENT)
			if err != nil {
				panic(err)
			}
			allRecords = append(allRecords, records...)
			bFound = found || found2
			blockStart -= 250
			blockEnd -= 250
		}
	} else { //沒有使用 -l 的話，就會把全部的reward event找完成
		totalBlocks := *totalB
		iterations := int(totalBlocks / 250)
		var remaining uint64 = uint64(totalBlocks % 250)
		for i := 0; i < iterations; i++ {
			records, err, _ := GetBlockEvents(ctx, blockStart, blockEnd, accessNodeClient, *allFlag, REWARD_PAID_EVENT)
			if err != nil {
				panic(err)
			}
			allRecords = append(allRecords, records...)
			records, err, _ = GetBlockEvents(ctx, blockStart, blockEnd, accessNodeClient, *allFlag, DELEGATOR_REWARD_PAID_EVENT)
			if err != nil {
				panic(err)
			}
			allRecords = append(allRecords, records...)
			blockStart -= 250
			blockEnd -= 250
		}
		if remaining != 0 {
			blockStart += uint64(250) - remaining
			records, err, _ := GetBlockEvents(ctx, blockStart, blockEnd, accessNodeClient, *allFlag, REWARD_PAID_EVENT)
			if err != nil {
				panic(err)
			}
			allRecords = append(allRecords, records...)
			records, err, _ = GetBlockEvents(ctx, blockStart, blockEnd, accessNodeClient, *allFlag, DELEGATOR_REWARD_PAID_EVENT)
			if err != nil {
				panic(err)
			}
			allRecords = append(allRecords, records...)
		}
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

// get RewardPaid and DelegatedRewardPaid
func GetBothRewardsBlockEvents(ctx context.Context, blockStart, blockEnd uint64, accessNodeClient *client.Client, all bool) ([]RewardRecord, error, bool) {
	// REWARD_PAID
	records1, err, bFound1 := GetBlockEvents(ctx, blockStart, blockEnd, accessNodeClient, all, REWARD_PAID_EVENT)
	if err != nil {
		return nil, err, false
	}
	// DELEGATOR_REWARD_PAID
	records2, err, bFound2 := GetBlockEvents(ctx, blockStart, blockEnd, accessNodeClient, all, DELEGATOR_REWARD_PAID_EVENT)
	if err != nil {
		return nil, err, false
	}

	// 合併兩個結果
	allRecords := append(records1, records2...)
	return allRecords, err, bFound1 || bFound2
}

// 用...搜尋進度，顯示搜尋結果
func GetBlockEvents(ctx context.Context, blockStart, blockEnd uint64, accessNodeClient *client.Client, all bool, eventString string) ([]RewardRecord, error, bool) {
	//	fmt.Println("Start = ", blockStart, "End = ", blockEnd);
	fmt.Print(".")
	// REWARD_PAID_ EVENT := "A.8624b52f9ddcd04a.FlowIDTableStaking.RewardsPaid"

	blockEvents, err := accessNodeClient.GetEventsForHeightRange(ctx,
		eventString,
		blockStart, blockEnd)
	/*
		51753836,
		51753836)
	*/
	bFound := false
	var results []RewardRecord

	if err != nil {
		panic(err)
		/*	    fmt.Println("Encounter errors:", err, "Block start =", blockStart, " block end = ", blockEnd)
			    return err, bFound
		*/
	}
	// twm node ID: e8f4bd649d08ecb5afb7023a0c5e8bb10ce56659399665da8cc9d517e7982f92

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
				// TWM_NODE_ID := "e8f4bd649d08ecb5afb7023a0c5e8bb10ce56659399665da8cc9d517e7982f92"
				if nodeID == TWM_NODE_ID {
					// 只保留 DelegatorID -1, 2, 3, 4
					if !(delegatorID == -1 || delegatorID == 2 || delegatorID == 3 || delegatorID == 4) {
						continue
					}
					// 獲取區塊時間戳
					blockByID, err := accessNodeClient.GetBlockByID(ctx, blockEvent.BlockID)
					if err != nil {
						fmt.Println("err:", err.Error())
						return nil, err, bFound
					}

					// 確定事件類型
					eventType := "Node"
					if strings.Contains(event.Type, "DelegatorRewardsPaid") {
						eventType = "Delegator"
					}

					// 存成 struct
					delegator2, delegator3, delegator4, err := GetDelegatorTotals()
					if err != nil {
						fmt.Println("查詢 delegator 總額失敗:", err)
						os.Exit(1)
					}
					nodeTotal, err := GetNodeTotal()
					if err != nil {
						fmt.Println("查詢 node_total 失敗:", err)
						os.Exit(1)
					}
					record := RewardRecord{
						Type:            eventType,
						NodeID:          nodeID,
						DelegatorID:     delegatorID,
						Amount:          amount,
						EpochCounter:    epochCounter,
						Timestamp:       blockByID.Timestamp.Format("2006-01-02 15:04:05"),
						DelegatorTotal2: delegator2,
						DelegatorTotal3: delegator3,
						DelegatorTotal4: delegator4,
						Node_total:      nodeTotal,
					}
					results = append(results, record)
					// 也 print 一下
					fmt.Printf("Type=%s|NodeID=%s|DelegatorID=%d|Amount=%.8f|EpochCounter=%d|Timestamp=%s\n",
						eventType,
						nodeID,
						delegatorID,
						amount,
						epochCounter,
						record.Timestamp)
					bFound = true
				}
				/*
				   				} else {
				       fmt.Print("X")
				   }
				*/
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

func printBlock(block *flow.Block, err error) {
	if err != nil {
		fmt.Println("err:", err.Error())
		panic(err)
	}
	fmt.Printf("\nID: %s\n", block.ID)
	fmt.Printf("height: %d\n", block.Height)
	fmt.Printf("timestamp: %s\n\n", block.Timestamp)
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

// 新增 function 查詢三個 delegator 的 tokensRewarded
func GetDelegatorTotals() (string, string, string, error) {
	url := "https://rest-mainnet.onflow.org/v1/scripts?blocks_height=final"
	method := "POST"
	ids := []int{2, 3, 4}
	totals := make(map[int]string)

	for _, delegatorID := range ids {
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
			return "", "", "", err
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return "", "", "", err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", "", "", err
		}

		cleaned := strings.Trim(string(body), "\"")
		decoded, err := base64.StdEncoding.DecodeString(cleaned)
		if err != nil {
			return "", "", "", err
		}

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
			return "", "", "", err
		}

		for _, field := range outer.Value.Fields {
			if field.Name == "tokensRewarded" {
				totals[delegatorID] = field.Value.Value
			}
		}
	}
	return totals[2], totals[3], totals[4], nil
}

// 新增 function 查詢 node_total (tokensRewarded)
func GetNodeTotal() (string, error) {
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
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	cleaned := strings.Trim(string(body), "\"")
	decoded, err := base64.StdEncoding.DecodeString(cleaned)
	if err != nil {
		return "", err
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
		return "", err
	}
	var nodeTotal string
	for _, field := range outer.Value.Fields {
		if field.Name == "tokensRewarded" {
			var v struct {
				Value string `json:"value"`
			}
			if err := json.Unmarshal(field.Value, &v); err == nil {
				nodeTotal = v.Value
			}
		}
	}
	return nodeTotal, nil
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
