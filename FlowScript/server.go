package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

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

	// 取得所有獎勵資料（依 timestamp DESC 排序），支援 ids 參數
	r.GET("/api/rewards/all", func(c *gin.Context) {
		db, err := NewDatabase("data.db")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "資料庫連線失敗"})
			return
		}
		defer db.Close()

		idsParam := c.Query("ids")
		var rewards []RewardData
		if idsParam != "" {
			idStrs := strings.Split(idsParam, ",")
			ids := []int{}
			for _, s := range idStrs {
				if s == "" {
					continue
				}
				id, err := strconv.Atoi(s)
				if err == nil {
					ids = append(ids, id)
				}
			}
			rewards, err = db.GetRewardsByDelegators(ids)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢失敗"})
				return
			}
		} else {
			rewards, err = db.GetAllRewards()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢失敗"})
				return
			}
		}
		c.JSON(http.StatusOK, rewards)
	})

	// 匯出 CSV，支援 ids 參數
	r.GET("/api/rewards/export", func(c *gin.Context) {
		db, err := NewDatabase("data.db")
		if err != nil {
			c.String(http.StatusInternalServerError, "資料庫連線失敗")
			return
		}
		defer db.Close()

		week := c.Query("week")
		idsParam := c.Query("ids")
		var rewards []RewardData
		if week == "all" || week == "" {
			if idsParam != "" {
				idStrs := strings.Split(idsParam, ",")
				ids := []int{}
				for _, s := range idStrs {
					if s == "" {
						continue
					}
					id, err := strconv.Atoi(s)
					if err == nil {
						ids = append(ids, id)
					}
				}
				rewards, err = db.GetRewardsByDelegators(ids)
				if err != nil {
					c.String(http.StatusInternalServerError, "查詢失敗")
					return
				}
			} else {
				rewards, err = db.GetAllRewards()
				if err != nil {
					c.String(http.StatusInternalServerError, "查詢失敗")
					return
				}
			}
		} else {
			parts := make([]string, 2)
			copy(parts, splitWeek(week))
			if parts[0] == "" || parts[1] == "" {
				c.String(http.StatusBadRequest, "週期參數錯誤")
				return
			}
			if idsParam != "" {
				idStrs := strings.Split(idsParam, ",")
				ids := []int{}
				for _, s := range idStrs {
					if s == "" {
						continue
					}
					id, err := strconv.Atoi(s)
					if err == nil {
						ids = append(ids, id)
					}
				}
				rewards, err = db.GetRewardsByDelegators(ids)
				if err != nil {
					c.String(http.StatusInternalServerError, "查詢失敗")
					return
				}
				// 篩選日期
				filtered := []RewardData{}
				for _, r := range rewards {
					if r.Timestamp.Format("2006-01-02") >= parts[0] && r.Timestamp.Format("2006-01-02") <= parts[1] {
						filtered = append(filtered, r)
					}
				}
				rewards = filtered
			} else {
				rewards, err = db.GetRewardsByDateRange(parts[0], parts[1])
				if err != nil {
					c.String(http.StatusInternalServerError, "查詢失敗")
					return
				}
			}
		}
		// 欄位順序: 獎勵發放時間,type,delegator_id,node_ID,epoch,本周Reward,delegator 可提領餘額,node 可提領餘額,delegator 質押數量,node 質押數量
		csv := "\uFEFF獎勵發放時間,type,delegator_id,node_ID,epoch,本周Reward,delegator 可提領餘額,node 可提領餘額,delegator 質押數量,node 質押數量\n" // UTF-8 BOM
		for _, row := range rewards {
			// 根據 delegator_id 決定 node 可提領餘額的值
			var nodeTotalValue float64
			if row.DelegatorID == -1 {
				nodeTotalValue = row.NodeTotal
			} else {
				nodeTotalValue = 0
			}

			// 根據 delegator_id 決定 node 質押數量的值
			var nodeStakedValue float64
			if row.DelegatorID == -1 {
				nodeStakedValue = row.NodeStaked
			} else {
				nodeStakedValue = 0
			}

			csv += "\"" + row.Timestamp.Format("2006-01-02 15:04:05") + "\"," +
				"\"" + row.Type + "\"," +
				"\"" + strconv.Itoa(row.DelegatorID) + "\"," +
				"\"FSD\"," +
				strconv.FormatUint(row.EpochCounter, 10) + "," +
				formatAmount(row.Amount) + "," +
				formatAmount(row.DelegatorTotal) + "," +
				formatAmount(nodeTotalValue) + "," +
				formatAmount(row.DelegatorStaked) + "," +
				formatAmount(nodeStakedValue) + "\n"
		}
		c.Header("Content-Type", "text/csv; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=reward_log.csv")
		c.String(http.StatusOK, csv)
	})

	// 啟動 server
	r.Run(":8081")
}

// 工具函數：分割 week 參數
func splitWeek(week string) []string {
	parts := make([]string, 2)
	arr := []rune(week)
	sep := -1
	for i, ch := range arr {
		if ch == '|' {
			sep = i
			break
		}
	}
	if sep != -1 {
		parts[0] = string(arr[:sep])
		parts[1] = string(arr[sep+1:])
	}
	return parts
}

// 工具函數：金額格式化 - 完整顯示小數點
func formatAmount(amount float64) string {
	return strconv.FormatFloat(amount, 'f', -1, 64)
}
