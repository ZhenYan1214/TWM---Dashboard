// db.go － SQLite 初始化、基礎 CRUD
package main /* ← 如果你要把它當獨立執行檔先測試，就保留 package main。
   之後要改成 library，只要把這行改成 package db 然後刪掉 main() 即可 */

import (
	"fmt"
	"log"
	"time"

	// ❗ go-sqlite3 使用 CGO，第一次執行前要先：go get github.com/mattn/go-sqlite3
	_ "github.com/mattn/go-sqlite3"

	"github.com/jmoiron/sqlx"
)

// Reward 跟 fetcher / models 定義相同，為了獨立可跑這裡再複製一份
type Reward struct {
	Epoch       uint64  `db:"epoch"`
	NodeID      string  `db:"node_id"`
	DelegatorID *uint32 `db:"delegator_id"`
	Amount      string  `db:"amount"`
	TsUTC       int64   `db:"ts_utc"`
}

// ---------- 公用函式（library 用） ----------

// Open 連線，若檔案不存在會自動建立；journal_mode=WAL 允許「讀取不中斷寫入」
func Open(path string) (*sqlx.DB, error) {
	return sqlx.Open("sqlite3",
		fmt.Sprintf("%s?_busy_timeout=10000&_journal_mode=WAL", path))
}

// Migrate 建表（如果已存在就跳過）
func Migrate(db *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS rewards (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		epoch        INTEGER NOT NULL,
		node_id      TEXT    NOT NULL,
		delegator_id INTEGER,
		amount       TEXT    NOT NULL,
		ts_utc       INTEGER NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_rewards_epoch ON rewards(epoch);
	`
	_, err := db.Exec(schema)
	return err
}

// InsertRewards 批次寫入
func InsertRewards(db *sqlx.DB, rows []Reward) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	for _, r := range rows {
		if _, err := tx.NamedExec(`
	INSERT INTO rewards(epoch,node_id,delegator_id,amount,ts_utc)
	VALUES(:epoch,:node_id,:delegator_id,:amount,:ts_utc)`, &r); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// ---------- 以下是「示範 main」，方便你直接 go run db.go 測試 ----------

func main() {
	// 1. 開 / 建 data.db
	db, err := Open("data.db")
	if err != nil {
		log.Fatalln("open db:", err)
	}
	defer db.Close()

	// 2. 建表
	if err := Migrate(db); err != nil {
		log.Fatalln("migrate:", err)
	}

	// 3. 示範塞一筆假資料（方便你開 DB Browser 看看）
	fake := Reward{
		Epoch:  999,
		NodeID: "demo-node-id",
		Amount: "123.4567",
		TsUTC:  time.Now().Unix(),
	}
	if err := InsertRewards(db, []Reward{fake}); err != nil {
		log.Fatalln("insert:", err)
	}

	log.Println("✅ data.db 已建立並插入 1 筆測試資料")
}
