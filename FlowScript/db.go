package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// RewardData 存儲獎勵資料的結構
// 支援所有 delegator，只有一個 delegator_total
type RewardData struct {
	Type            string    `json:"type"`
	NodeID          string    `json:"node_id"`
	DelegatorID     int       `json:"delegator_id"`
	Amount          float64   `json:"amount"`
	EpochCounter    uint64    `json:"epoch_counter"`
	Timestamp       time.Time `json:"timestamp"`
	DelegatorTotal  float64   `json:"delegator_total"`
	NodeTotal       float64   `json:"node_total"`
	NodeStaked      float64   `json:"node_staked"`
	DelegatorStaked float64   `json:"delegator_staked"`
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
		delegator_total REAL NOT NULL,
		node_total REAL NOT NULL,
		node_staked REAL NOT NULL,
		delegator_staked REAL NOT NULL
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
	INSERT INTO rewards (type, node_id, delegator_id, amount, epoch_counter, timestamp, delegator_total, node_total, node_staked, delegator_staked)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		reward.Type,
		reward.NodeID,
		reward.DelegatorID,
		reward.Amount,
		reward.EpochCounter,
		reward.Timestamp,
		reward.DelegatorTotal,
		reward.NodeTotal,
		reward.NodeStaked,
		reward.DelegatorStaked,
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
		INSERT INTO rewards (type, node_id, delegator_id, amount, epoch_counter, timestamp, delegator_total, node_total, node_staked, delegator_staked)
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
			reward.DelegatorTotal,
			reward.NodeTotal,
			reward.NodeStaked,
			reward.DelegatorStaked,
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
	SELECT type, node_id, delegator_id, amount, epoch_counter, timestamp, delegator_total, node_total, node_staked, delegator_staked
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
			&reward.DelegatorTotal,
			&reward.NodeTotal,
			&reward.NodeStaked,
			&reward.DelegatorStaked,
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
	SELECT type, node_id, delegator_id, amount, epoch_counter, timestamp, delegator_total, node_total, node_staked, delegator_staked
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
			&reward.DelegatorTotal,
			&reward.NodeTotal,
			&reward.NodeStaked,
			&reward.DelegatorStaked,
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
	SELECT type, node_id, delegator_id, amount, epoch_counter, timestamp, delegator_total, node_total, node_staked, delegator_staked
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
			&reward.DelegatorTotal,
			&reward.NodeTotal,
			&reward.NodeStaked,
			&reward.DelegatorStaked,
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
        SELECT type, node_id, delegator_id, amount, epoch_counter, timestamp, delegator_total, node_total, node_staked, delegator_staked
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
			&reward.DelegatorTotal,
			&reward.NodeTotal,
			&reward.NodeStaked,
			&reward.DelegatorStaked,
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

// 根據日期範圍查詢獎勵記錄，用來匯出CSV
func (d *Database) GetRewardsByDateRange(start, end string) ([]RewardData, error) {
	query := `
	SELECT type, node_id, delegator_id, amount, epoch_counter, timestamp, delegator_total, node_total, node_staked, delegator_staked
	FROM rewards
	WHERE timestamp >= ? AND timestamp <= ?
	ORDER BY timestamp DESC
	`
	rows, err := d.db.Query(query, start+" 00:00:00", end+" 23:59:59")
	if err != nil {
		return nil, fmt.Errorf("failed to query rewards by date range: %w", err)
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
			&reward.DelegatorTotal,
			&reward.NodeTotal,
			&reward.NodeStaked,
			&reward.DelegatorStaked,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reward: %w", err)
		}
		rewards = append(rewards, reward)
	}
	return rewards, nil
}

// 根據多個 delegatorId 查詢獎勵記錄
func (d *Database) GetRewardsByDelegators(ids []int) ([]RewardData, error) {
	if len(ids) == 0 {
		return d.GetAllRewards()
	}
	query := `SELECT type, node_id, delegator_id, amount, epoch_counter, timestamp, delegator_total, node_total, node_staked, delegator_staked FROM rewards WHERE delegator_id IN (`
	params := make([]interface{}, len(ids))
	for i, id := range ids {
		params[i] = id
	}
	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
	}
	query += strings.Join(placeholders, ",") + ") ORDER BY timestamp DESC"
	rows, err := d.db.Query(query, params...)
	if err != nil {
		return nil, err
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
			&reward.DelegatorTotal,
			&reward.NodeTotal,
			&reward.NodeStaked,
			&reward.DelegatorStaked,
		)
		if err != nil {
			return nil, err
		}
		rewards = append(rewards, reward)
	}
	return rewards, nil
}
