package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	DB_PATH = "data.db"
)

// Application 主應用程式結構
type Application struct {
	db      *Database
	fetcher *RewardFetcher
}

// NewApplication 建立新的應用程式實例
func NewApplication() (*Application, error) {
	// 初始化數據庫
	db, err := NewDatabase(DB_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// 初始化獎勵抓取器
	fetcher, err := NewRewardFetcher()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize fetcher: %w", err)
	}

	return &Application{
		db:      db,
		fetcher: fetcher,
	}, nil
}

// fetchCommand 執行抓取命令
func (app *Application) fetchCommand() error {
	fmt.Println("🚀 Starting to fetch rewards...")

	startTime := time.Now()

	// 從 Flow 區塊鏈抓取獎勵數據
	rewards, err := app.fetcher.FetchRewards()
	if err != nil {
		return fmt.Errorf("failed to fetch rewards: %w", err)
	}

	if len(rewards) == 0 {
		fmt.Println("⚠️  No rewards found")
		return nil
	}

	// 將數據存入數據庫
	fmt.Printf("💾 Saving %d rewards to database...\n", len(rewards))
	err = app.db.InsertRewards(rewards)
	if err != nil {
		return fmt.Errorf("failed to save rewards: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("✅ Successfully fetched and saved %d rewards in %v\n", len(rewards), duration)

	// 顯示抓取的數據
	app.displayRewards(rewards)

	return nil
}

// listCommand 列出數據庫中的所有獎勵
func (app *Application) listCommand() error {
	fmt.Println("📋 Listing all rewards from database...")

	rewards, err := app.db.GetAllRewards()
	if err != nil {
		return fmt.Errorf("failed to get rewards: %w", err)
	}

	if len(rewards) == 0 {
		fmt.Println("⚠️  No rewards found in database")
		return nil
	}

	app.displayRewards(rewards)
	return nil
}

// statsCommand 顯示獎勵統計
func (app *Application) statsCommand() error {
	fmt.Println("📊 Displaying reward statistics...")

	stats, err := app.db.GetRewardStats()
	if err != nil {
		return fmt.Errorf("failed to get stats: %w", err)
	}

	fmt.Printf("\n=== Reward Statistics ===\n")
	fmt.Printf("Total Amount: %.8f FLOW\n", stats["total_amount"])
	fmt.Printf("Total Records: %d\n", stats["total_count"])

	fmt.Println("\nBy Type:")
	typeStats := stats["type_stats"].(map[string]int)
	for rewardType, count := range typeStats {
		fmt.Printf("  %s: %d records\n", rewardType, count)
	}

	return nil
}

// clearCommand 清空數據庫
func (app *Application) clearCommand() error {
	fmt.Print("⚠️  Are you sure you want to clear all rewards? (y/N): ")
	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" {
		fmt.Println("Operation cancelled")
		return nil
	}

	err := app.db.ClearRewards()
	if err != nil {
		return fmt.Errorf("failed to clear rewards: %w", err)
	}

	fmt.Println("✅ All rewards cleared successfully")
	return nil
}

// displayRewards 顯示獎勵列表
func (app *Application) displayRewards(rewards []RewardData) {
	fmt.Printf("\n=== Rewards (%d records) ===\n", len(rewards))
	fmt.Printf("%-20s %-12s %-10s %-15s %-12s %-20s\n",
		"Type", "DelegatorID", "Amount", "EpochCounter", "BlockHeight", "Timestamp")
	fmt.Println(strings.Repeat("-", 100))

	for _, reward := range rewards {
		delegatorStr := fmt.Sprintf("%d", reward.DelegatorID)
		if reward.DelegatorID == -1 {
			delegatorStr = "N/A"
		}

		fmt.Printf("%-20s %-12s %-10.8f %-15d %-12d %-20s\n",
			reward.Type,
			delegatorStr,
			reward.Amount,
			reward.EpochCounter,
			reward.BlockHeight,
			reward.Timestamp.Format("2006-01-02 15:04:05"))
	}
}

// showUsage 顯示使用說明
func (app *Application) showUsage() {
	fmt.Println(`
Flow Rewards CLI Tool

Usage:
  go run CLI.go [command]

Commands:
  fetch, f     Fetch rewards from Flow blockchain and save to database
  list, l      List all rewards from database
  stats, s     Show reward statistics
  clear, c     Clear all rewards from database
  help, h      Show this help message

Examples:
  go run CLI.go fetch     # Fetch new rewards
  go run CLI.go list      # List all rewards
  go run CLI.go stats     # Show statistics
`)
}

// Close 關閉資源
func (app *Application) Close() {
	if app.db != nil {
		app.db.Close()
	}
}

func main() {
	app, err := NewApplication()
	if err != nil {
		fmt.Printf("❌ Failed to initialize application: %v\n", err)
		os.Exit(1)
	}
	defer app.Close()

	// 處理命令行參數
	if len(os.Args) < 2 {
		app.showUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "fetch", "f":
		err = app.fetchCommand()
	case "list", "l":
		err = app.listCommand()
	case "stats", "s":
		err = app.statsCommand()
	case "clear", "c":
		err = app.clearCommand()
	case "help", "h":
		app.showUsage()
		return
	default:
		fmt.Printf("❌ Unknown command: %s\n", command)
		app.showUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
}
