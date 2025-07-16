package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// RewardRecord struct 必須與 fetch.go 相同
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

// RewardData 存儲獎勵資料的結構
type RewardData struct {
	Type             string    `json:"type"`
	NodeID           string    `json:"node_id"`
	DelegatorID      int       `json:"delegator_id"`
	Amount           float64   `json:"amount"`
	EpochCounter     uint64    `json:"epoch_counter"`
	Timestamp        time.Time `json:"timestamp"`
	Delegator_total2 float64   `json:"delegator_total2"`
	Delegator_total3 float64   `json:"delegator_total3"`
	Delegator_total4 float64   `json:"delegator_total4"`
	Node_total       float64   `json:"node_total"`
}

// Database 封裝數據庫操作
type Database struct {
	db *sql.DB
}

// NewDatabase 建立新的數據庫連接
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	database := &Database{db: db}

	// 建立資料表
	if err := database.createTable(); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return database, nil
}

// createTable 建立 rewards 資料表
func (d *Database) createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS rewards (
		type TEXT NOT NULL,
		node_id TEXT NOT NULL,
		delegator_id INTEGER NOT NULL,
		amount REAL NOT NULL,
		epoch_counter INTEGER NOT NULL,
		timestamp DATETIME NOT NULL,
		delegator_total2 REAL DEFAULT 0,
		delegator_total3 REAL DEFAULT 0,
		delegator_total4 REAL DEFAULT 0,
		node_total REAL DEFAULT 0
	);
	
	CREATE INDEX IF NOT EXISTS idx_type ON rewards(type);
	CREATE INDEX IF NOT EXISTS idx_delegator_id ON rewards(delegator_id);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON rewards(timestamp);
	`

	_, err := d.db.Exec(query)
	return err
}

// InsertRewards 批量插入獎勵記錄
func (d *Database) InsertRewards(rewards []RewardData) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO rewards (type, node_id, delegator_id, amount, epoch_counter, timestamp, delegator_total2, delegator_total3, delegator_total4, node_total)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, reward := range rewards {
		_, err = stmt.Exec(
			reward.Type,
			reward.NodeID,
			reward.DelegatorID,
			reward.Amount,
			reward.EpochCounter,
			reward.Timestamp,
			reward.Delegator_total2,
			reward.Delegator_total3,
			reward.Delegator_total4,
			reward.Node_total,
		)
		if err != nil {
			return fmt.Errorf("failed to insert reward: %w", err)
		}
	}

	return tx.Commit()
}

// Close 關閉數據庫連接
func (d *Database) Close() error {
	return d.db.Close()
}

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
	data, err := ioutil.ReadFile("rewards.json")
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

	// 轉換 RewardRecord -> RewardData
	var rewardDataList []RewardData
	for _, r := range rewards {
		t, _ := time.Parse("2006-01-02 15:04:05", r.Timestamp)
		var d2, d3, d4, nodeTotal float64
		fmt.Sscanf(r.DelegatorTotal2, "%f", &d2)
		fmt.Sscanf(r.DelegatorTotal3, "%f", &d3)
		fmt.Sscanf(r.DelegatorTotal4, "%f", &d4)
		fmt.Sscanf(r.Node_total, "%f", &nodeTotal)
		rewardDataList = append(rewardDataList, RewardData{
			Type:             r.Type,
			NodeID:           r.NodeID,
			DelegatorID:      r.DelegatorID,
			Amount:           r.Amount,
			EpochCounter:     r.EpochCounter,
			Timestamp:        t,
			Delegator_total2: d2,
			Delegator_total3: d3,
			Delegator_total4: d4,
			Node_total:       nodeTotal,
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
