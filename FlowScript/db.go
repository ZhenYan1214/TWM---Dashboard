package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

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

// InsertReward 插入單一獎勵記錄
func (d *Database) InsertReward(reward RewardData) error {
	query := `
	INSERT INTO rewards (type, node_id, delegator_id, amount, epoch_counter, timestamp, delegator_total2, delegator_total3, delegator_total4, node_total)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
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

	return nil
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

// GetAllRewards 取得所有獎勵記錄
func (d *Database) GetAllRewards() ([]RewardData, error) {
	query := `
	SELECT type, node_id, delegator_id, amount, epoch_counter, timestamp, delegator_total2, delegator_total3, delegator_total4, node_total
	FROM rewards
	ORDER BY timestamp DESC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query rewards: %w", err)
	}
	defer rows.Close()

	var rewards []RewardData
	for rows.Next() {
		var reward RewardData
		err := rows.Scan(
			&reward.Type,
			&reward.NodeID,
			&reward.DelegatorID,
			&reward.Amount,
			&reward.EpochCounter,
			&reward.Timestamp,
			&reward.Delegator_total2,
			&reward.Delegator_total3,
			&reward.Delegator_total4,
			&reward.Node_total,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reward: %w", err)
		}
		rewards = append(rewards, reward)
	}

	return rewards, nil
}

// GetRewardsByType 根據類型取得獎勵記錄
func (d *Database) GetRewardsByType(rewardType string) ([]RewardData, error) {
	query := `
	SELECT type, node_id, delegator_id, amount, epoch_counter, timestamp, delegator_total2, delegator_total3, delegator_total4, node_total
	FROM rewards
	WHERE type = ?
	ORDER BY timestamp DESC
	`

	rows, err := d.db.Query(query, rewardType)
	if err != nil {
		return nil, fmt.Errorf("failed to query rewards by type: %w", err)
	}
	defer rows.Close()

	var rewards []RewardData
	for rows.Next() {
		var reward RewardData
		err := rows.Scan(
			&reward.Type,
			&reward.NodeID,
			&reward.DelegatorID,
			&reward.Amount,
			&reward.EpochCounter,
			&reward.Timestamp,
			&reward.Delegator_total2,
			&reward.Delegator_total3,
			&reward.Delegator_total4,
			&reward.Node_total,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reward: %w", err)
		}
		rewards = append(rewards, reward)
	}

	return rewards, nil
}

// GetRewardsByDelegator 根據委託者 ID 取得獎勵記錄
func (d *Database) GetRewardsByDelegator(delegatorID int) ([]RewardData, error) {
	query := `
	SELECT type, node_id, delegator_id, amount, epoch_counter, timestamp, delegator_total2, delegator_total3, delegator_total4, node_total
	FROM rewards
	WHERE delegator_id = ?
	ORDER BY timestamp DESC
	`

	rows, err := d.db.Query(query, delegatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to query rewards by delegator: %w", err)
	}
	defer rows.Close()

	var rewards []RewardData
	for rows.Next() {
		var reward RewardData
		err := rows.Scan(
			&reward.Type,
			&reward.NodeID,
			&reward.DelegatorID,
			&reward.Amount,
			&reward.EpochCounter,
			&reward.Timestamp,
			&reward.Delegator_total2,
			&reward.Delegator_total3,
			&reward.Delegator_total4,
			&reward.Node_total,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reward: %w", err)
		}
		rewards = append(rewards, reward)
	}

	return rewards, nil
}

// GetRewardStats 取得獎勵統計資訊
func (d *Database) GetRewardStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 總獎勵金額
	var totalAmount float64
	err := d.db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM rewards").Scan(&totalAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total amount: %w", err)
	}
	stats["total_amount"] = totalAmount

	// 總記錄數
	var totalCount int
	err = d.db.QueryRow("SELECT COUNT(*) FROM rewards").Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}
	stats["total_count"] = totalCount

	// 各類型統計
	typeStats := make(map[string]int)
	rows, err := d.db.Query("SELECT type, COUNT(*) FROM rewards GROUP BY type")
	if err != nil {
		return nil, fmt.Errorf("failed to get type stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var rewardType string
		var count int
		if err := rows.Scan(&rewardType, &count); err != nil {
			return nil, fmt.Errorf("failed to scan type stats: %w", err)
		}
		typeStats[rewardType] = count
	}
	stats["type_stats"] = typeStats

	return stats, nil
}

// 取得最新一批（timestamp 最大）獎勵資料
func (d *Database) GetLatestBatchRewards() ([]RewardData, error) {
	// 查詢最新的 timestamp
	var latestTimestamp string
	err := d.db.QueryRow(`SELECT MAX(timestamp) FROM rewards`).Scan(&latestTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest timestamp: %w", err)
	}

	// 查詢該 timestamp 的所有資料
	query := `
        SELECT type, node_id, delegator_id, amount, epoch_counter, timestamp, delegator_total2, delegator_total3, delegator_total4, node_total
        FROM rewards
        WHERE timestamp = ?
        ORDER BY node_id
    `
	rows, err := d.db.Query(query, latestTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to query latest batch rewards: %w", err)
	}
	defer rows.Close()

	var rewards []RewardData
	for rows.Next() {
		var reward RewardData
		err := rows.Scan(
			&reward.Type,
			&reward.NodeID,
			&reward.DelegatorID,
			&reward.Amount,
			&reward.EpochCounter,
			&reward.Timestamp,
			&reward.Delegator_total2,
			&reward.Delegator_total3,
			&reward.Delegator_total4,
			&reward.Node_total,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reward: %w", err)
		}
		rewards = append(rewards, reward)
	}
	return rewards, nil
}

// Close 關閉數據庫連接
func (d *Database) Close() error {
	return d.db.Close()
}

// ClearRewards 清空所有獎勵記錄（用於測試）
func (d *Database) ClearRewards() error {
	_, err := d.db.Exec("DELETE FROM rewards")
	return err
}

// main 函式已移除，因為 CLI.go 會使用這個檔案
