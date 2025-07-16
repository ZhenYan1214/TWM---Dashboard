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
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/onflow/flow-go-sdk"
	client "github.com/onflow/flow-go-sdk/access/grpc"

	// required for increasing the max gRPC size
	// required for increasing the max gRPC size
	"github.com/onflow/cadence"
)

// twm node ID: e8f4bd649d08ecb5afb7023a0c5e8bb10ce56659399665da8cc9d517e7982f92
const TWM_NODE_ID = "e8f4bd649d08ecb5afb7023a0c5e8bb10ce56659399665da8cc9d517e7982f92"
const REWARD_PAID_EVENT = "A.8624b52f9ddcd04a.FlowIDTableStaking.RewardsPaid"
const DELEGATOR_REWARD_PAID_EVENT = "A.8624b52f9ddcd04a.FlowIDTableStaking.DelegatorRewardsPaid"

// increase max gRPC message size to 16 MB to resolve the issue for DelegatedRewardsPaid, default message size is 4 MB
const MAX_GRPC_MESSAGE_SIZE = 1024 * 1024 * 16

// An Identifier is a 32-byte unique identifier for an entity. From flow.go
type Identifier [32]byte

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
	// the Flow testnet community Access node API endpoint
	accessNodeAddress := "access.mainnet.nodes.onflow.org:9000"

	// create a gRPC client for the Access node
	// accessNodeClient, err := client.NewClient(accessNodeAddress)
	// change the default gRPC message size
	accessNodeClient, err := client.NewClient(accessNodeAddress)
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

	// check if we want to jump start to Wed night 20:00
	if *jumpStart {
		fmt.Println("Jump to start height on Wednesday ...")
		blockEnd, err = JumpToNewEndHeight(ctx, accessNodeClient, blockEnd)
		if err != nil {
			panic(err)
		}
		fmt.Println("Find new end height =", blockEnd)
	}
	bFound := false
	blockStart := blockEnd - 249
	if *lastReward {
		// loop until find the last reward information
		for bFound == false {
			err, bFound = GetBothRewardsBlockEvents(ctx, blockStart, blockEnd, accessNodeClient, *allFlag)
			if err != nil {
				panic(err)
			}
			blockStart -= 250
			blockEnd -= 250
		}
	} else {
		// looping for more blocks
		totalBlocks := *totalB
		iterations := int(totalBlocks / 250)
		var remaining uint64 = uint64(totalBlocks % 250)
		for i := 0; i < iterations; i++ {
			err, bFound = GetBothRewardsBlockEvents(ctx, blockStart, blockEnd, accessNodeClient, *allFlag)
			if err != nil {
				panic(err)
			}
			blockStart -= 250
			blockEnd -= 250
		}
		if remaining != 0 {
			blockStart += uint64(250) - remaining
			// REWARD_PAID and DELEGATOR_REWARD_PAID
			err, bFound = GetBothRewardsBlockEvents(ctx, blockStart, blockEnd, accessNodeClient, *allFlag)
			if err != nil {
				panic(err)
			}
		}
	}
	// print the finished time
	finishedTime := time.Now()
	// format the time stamp as "YYYY-MM-DD HH:MM:SS"
	timeStamp = finishedTime.Format("2006-01-02 15:04:05")
	fmt.Println("Finished time:", timeStamp)
	fmt.Println("***")
}

// get RewardPaid and DelegatedRewardPaid
func GetBothRewardsBlockEvents(ctx context.Context, blockStart, blockEnd uint64, accessNodeClient *client.Client, all bool) (error, bool) {
	// REWARD_PAID
	err, bFound := GetBlockEvents(ctx, blockStart, blockEnd, accessNodeClient, all, REWARD_PAID_EVENT)
	if err != nil {
		panic(err)
	}
	// DELEGATOR_REWARD_PAID
	err, bFound = GetBlockEvents(ctx, blockStart, blockEnd, accessNodeClient, all, DELEGATOR_REWARD_PAID_EVENT)
	if err != nil {
		panic(err)
	}
	return err, bFound
}

func GetBlockEvents(ctx context.Context, blockStart, blockEnd uint64, accessNodeClient *client.Client, all bool, eventString string) (error, bool) {
	fmt.Print(".")

	blockEvents, err := accessNodeClient.GetEventsForHeightRange(ctx,
		eventString,
		blockStart, blockEnd)
	bFound := false

	if err != nil {
		panic(err)
	}

	type SimpleEvent struct {
		EventType    string `json:"eventType"`
		NodeID       string `json:"nodeID"`
		DelegatorID  *int   `json:"delegatorID,omitempty"`
		Amount       string `json:"amount"`
		EpochCounter int    `json:"epochCounter"`
		Timestamp    string `json:"timestamp"`
	}

	for _, blockEvent := range blockEvents {
		// 取得 block timestamp
		blockByID, err := accessNodeClient.GetBlockByID(ctx, blockEvent.BlockID)
		var timestamp string
		if err == nil {
			timestamp = blockByID.Timestamp.Format(time.RFC3339)
		} else {
			timestamp = ""
		}
		for _, event := range blockEvent.Events {
			if all {
				fmt.Println("\tEvent type: " + event.Type)
				fmt.Println("\tEvent: " + event.Value.String())
				fmt.Println("\tEvent payload: " + string(event.Payload))
			} else {
				fm := event.Value.FieldsMappedByName()
				var nodeID string
				var delegatorID *int
				var amount string
				var epochCounter int
				if event.Type == REWARD_PAID_EVENT {
					if fm["id"] != nil && fm["amount"] != nil && fm["epochCounter"] != nil {
						nodeID = string(fm["id"].(cadence.String))
						amount = fm["amount"].(cadence.UFix64).String()
						epochCounter = int(uint64(fm["epochCounter"].(cadence.UInt64)))
					}
				} else if event.Type == DELEGATOR_REWARD_PAID_EVENT {
					if fm["id"] != nil && fm["delegatorID"] != nil && fm["amount"] != nil && fm["epochCounter"] != nil {
						nodeID = string(fm["id"].(cadence.String))
						id := int(uint32(fm["delegatorID"].(cadence.UInt32)))
						delegatorID = &id
						amount = fm["amount"].(cadence.UFix64).String()
						epochCounter = int(uint64(fm["epochCounter"].(cadence.UInt64)))
					}
				}
				if nodeID == TWM_NODE_ID {
					eventTypeShort := event.Type[strings.LastIndex(event.Type, ".")+1:]
					simpleEvent := SimpleEvent{
						EventType:    eventTypeShort,
						NodeID:       nodeID,
						DelegatorID:  delegatorID,
						Amount:       amount,
						EpochCounter: epochCounter,
						Timestamp:    timestamp,
					}
					jsonBytes, _ := json.Marshal(simpleEvent)
					fmt.Println(string(jsonBytes))
					bFound = true
				}
			}
		}
	}
	return err, bFound
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

// extract node ID, delegatorID, and amount from the event
func extractNodeIDAndAmount(eventString string) (string, int, float64, error) {
	// Regular expression pattern to match node ID and amount
	// pattern := `nodeID: "([^"]+)", amount: ([0-9.]+)`
	// revised to extract delegator ID information as well
	pattern := `nodeID: "([^"]+)", (?:\s*delegatorID: (\d+),)?\s*amount: ([0-9.]+)`

	// Compile the regular expression pattern
	re := regexp.MustCompile(pattern)

	// Find the matches in the event string
	matches := re.FindStringSubmatch(eventString)

	lenmatches := len(matches)

	if lenmatches < 3 {
		return "", -1, 0, fmt.Errorf("unable to extract node ID and amount from event string")
	}

	nodeID := strings.TrimSpace(matches[1])
	var delegatorID int
	if lenmatches == 4 {
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

/* sample reward paid message
	Event type: A.8624b52f9ddcd04a.FlowIDTableStaking.RewardsPaid
	Event: A.8624b52f9ddcd04a.FlowIDTableStaking.RewardsPaid(nodeID: "5a4bff17941a73909472afe23f1ccdc59d7526f93b16b4e374bd8353f8b624b4", amount: 9223.47334447)
	Event payload: {"value":{"id":"A.8624b52f9ddcd04a.FlowIDTableStaking.RewardsPaid","fields":[{"value":{"value":"5a4bff17941a73909472afe23f1ccdc59d7526f93b16b4e374bd8353f8b624b4","type":"String"},"name":"nodeID"},{"value":{"value":"9223.47334447","type":"UFix64"},"name":"amount"}]},"type":"Event"}

sample delegator reward paid message
	Event type: A.8624b52f9ddcd04a.FlowIDTableStaking.DelegatorRewardsPaid
	Event: A.8624b52f9ddcd04a.FlowIDTableStaking.DelegatorRewardsPaid(nodeID: "95ffacf0c05757cff71a4ee49e025d5a6d1103a3aa7d91253079e1bfb7c22458", delegatorID: 23, amount: 0.10424555)
	Event payload: {"value":{"id":"A.8624b52f9ddcd04a.FlowIDTableStaking.DelegatorRewardsPaid","fields":[{"value":{"value":"95ffacf0c05757cff71a4ee49e025d5a6d1103a3aa7d91253079e1bfb7c22458","type":"String"},"name":"nodeID"},{"value":{"value":"23","type":"UInt32"},"name":"delegatorID"},{"value":{"value":"0.10424555","type":"UFix64"},"name":"amount"}]},"type":"Event"}

*/
