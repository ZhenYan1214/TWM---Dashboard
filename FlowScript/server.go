package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 取得最新一批獎勵資料（timestamp 最大的那批）
	r.GET("/api/rewards", func(c *gin.Context) {
		db, err := NewDatabase("data.db")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "資料庫連線失敗"})
			return
		}
		defer db.Close()

		rewards, err := db.GetLatestBatchRewards()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢失敗"})
			return
		}
		c.JSON(http.StatusOK, rewards)
	})

	// 依類型查詢，後續會做一個根據type來顯示內容
	r.GET("/api/rewards/type/:type", func(c *gin.Context) {
		db, err := NewDatabase("data.db")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "資料庫連線失敗"})
			return
		}
		defer db.Close()

		rewardType := c.Param("type")
		rewards, err := db.GetRewardsByType(rewardType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢失敗"})
			return
		}
		c.JSON(http.StatusOK, rewards)
	})

	// 依 delegator_id 查詢，可以製作歷史查詢
	r.GET("/api/rewards/delegator/:id", func(c *gin.Context) {
		db, err := NewDatabase("data.db")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "資料庫連線失敗"})
			return
		}
		defer db.Close()

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "參數錯誤"})
			return
		}
		rewards, err := db.GetRewardsByDelegator(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢失敗"})
			return
		}
		c.JSON(http.StatusOK, rewards)
	})

	// 啟動 server
	r.Run(":8080")
}
