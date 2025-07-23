package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// RewardRecord struct 必須與 fetch.go 產生的 reward.json 欄位一致
// 只要 main.go 能正確解析 reward.json 即可
type RewardRecord struct {
	Type           string  `json:"type"`
	NodeID         string  `json:"node_id"`
	DelegatorID    int     `json:"delegator_id"`
	Amount         float64 `json:"amount"`
	EpochCounter   uint64  `json:"epoch_counter"`
	Timestamp      string  `json:"timestamp"`
	DelegatorTotal string  `json:"delegator_total"`
	NodeTotal      string  `json:"node_total"`
}

// 移除 RewardData struct
// 移除 Database struct
// 移除 NewDatabase func

func main() {
	// 1. 執行 fetch.go
	cmd := exec.Command("go", "run", "fetch.go", "-l", "-j")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("fetch.go 執行失敗: %v\n", err)
		os.Exit(1)
	}

	// 2. 讀取 rewards.json
	data, err := os.ReadFile("rewards.json")
	if err != nil {
		fmt.Printf("讀取 rewards.json 失敗: %v\n", err)
		os.Exit(1)
	}
	var rewards []RewardRecord
	if err := json.Unmarshal(data, &rewards); err != nil {
		fmt.Printf("解析 rewards.json 失敗: %v\n", err)
		os.Exit(1)
	}

	if len(rewards) == 0 {
		fmt.Println("No rewards found")
		return
	}

	// 3. 存進資料庫
	db, err := NewDatabase("data.db")
	if err != nil {
		fmt.Printf("資料庫初始化失敗: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// 轉換 RewardRecord -> RewardData（每個 delegator 都存一筆）
	var rewardDataList []RewardData
	for _, r := range rewards {
		t, _ := time.Parse("2006-01-02 15:04:05", r.Timestamp)
		delegatorTotal := 0.0
		fmt.Sscanf(r.DelegatorTotal, "%f", &delegatorTotal)
		nodeTotal := 0.0
		fmt.Sscanf(r.NodeTotal, "%f", &nodeTotal)
		rewardDataList = append(rewardDataList, RewardData{
			Type:           r.Type,
			NodeID:         r.NodeID,
			DelegatorID:    r.DelegatorID,
			Amount:         r.Amount,
			EpochCounter:   r.EpochCounter,
			Timestamp:      t,
			DelegatorTotal: delegatorTotal,
			NodeTotal:      nodeTotal,
		})
	}

	err = db.InsertRewards(rewardDataList)
	if err != nil {
		fmt.Printf("存入資料庫失敗: %v\n", err)
		os.Exit(1)
	}

	// 4. 刪除 rewards.json
	if err := os.Remove("rewards.json"); err != nil {
		fmt.Println("⚠️ 刪除 rewards.json 失敗，但不影響主流程：", err)
	}

	fmt.Printf("✅ Successfully fetched and saved %d rewards!\n", len(rewardDataList))
}
