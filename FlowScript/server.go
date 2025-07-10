package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 取得所有獎勵資料
	r.GET("/api/rewards", func(c *gin.Context) {
		db, err := NewDatabase("data.db")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "資料庫連線失敗"})
			return
		}
		defer db.Close()

		rewards, err := db.GetAllRewards()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢失敗"})
			return
		}
		c.JSON(http.StatusOK, rewards)
	})

	// 依類型查詢
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

	// 依 delegator_id 查詢
	r.GET("/api/rewards/delegator/:id", func(c *gin.Context) {
		db, err := NewDatabase("data.db")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "資料庫連線失敗"})
			return
		}
		defer db.Close()

		var id int
		if err := c.ShouldBindUri(&id); err != nil {
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
