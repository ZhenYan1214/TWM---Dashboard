package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	API_PORT = "8080"
	DB_PATH  = "data.db"
)

// APIServer 封裝 API 服務器
type APIServer struct {
	db     *Database
	router *mux.Router
}

// APIResponse 統一的 API 響應格式
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// NewAPIServer 建立新的 API 服務器
func NewAPIServer() (*APIServer, error) {
	db, err := NewDatabase(DB_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	server := &APIServer{
		db:     db,
		router: mux.NewRouter(),
	}

	server.setupRoutes()
	return server, nil
}

// setupRoutes 設置 API 路由
func (s *APIServer) setupRoutes() {
	api := s.router.PathPrefix("/api").Subrouter()

	// 獎勵相關路由
	api.HandleFunc("/rewards", s.getAllRewards).Methods("GET")
	api.HandleFunc("/rewards/type/{type}", s.getRewardsByType).Methods("GET")
	api.HandleFunc("/rewards/delegator/{id}", s.getRewardsByDelegator).Methods("GET")
	api.HandleFunc("/rewards/stats", s.getRewardStats).Methods("GET")

	// 操作相關路由
	api.HandleFunc("/fetch", s.fetchRewards).Methods("POST")
	api.HandleFunc("/clear", s.clearRewards).Methods("DELETE")

	// 健康檢查
	api.HandleFunc("/health", s.healthCheck).Methods("GET")

	// 靜態文件服務（如果需要）
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
}

// getAllRewards 取得所有獎勵
func (s *APIServer) getAllRewards(w http.ResponseWriter, r *http.Request) {
	rewards, err := s.db.GetAllRewards()
	if err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, "Failed to get rewards", err)
		return
	}

	s.sendSuccessResponse(w, "Rewards retrieved successfully", rewards)
}

// getRewardsByType 根據類型取得獎勵
func (s *APIServer) getRewardsByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rewardType := vars["type"]

	// 驗證類型
	if rewardType != "RewardsPaid" && rewardType != "DelegatorRewardsPaid" {
		s.sendErrorResponse(w, http.StatusBadRequest, "Invalid reward type", nil)
		return
	}

	rewards, err := s.db.GetRewardsByType(rewardType)
	if err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, "Failed to get rewards by type", err)
		return
	}

	s.sendSuccessResponse(w, fmt.Sprintf("Rewards of type %s retrieved successfully", rewardType), rewards)
}

// getRewardsByDelegator 根據委託者 ID 取得獎勵
func (s *APIServer) getRewardsByDelegator(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	delegatorIDStr := vars["id"]

	delegatorID, err := strconv.Atoi(delegatorIDStr)
	if err != nil {
		s.sendErrorResponse(w, http.StatusBadRequest, "Invalid delegator ID", err)
		return
	}

	rewards, err := s.db.GetRewardsByDelegator(delegatorID)
	if err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, "Failed to get rewards by delegator", err)
		return
	}

	s.sendSuccessResponse(w, fmt.Sprintf("Rewards for delegator %d retrieved successfully", delegatorID), rewards)
}

// getRewardStats 取得獎勵統計
func (s *APIServer) getRewardStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.db.GetRewardStats()
	if err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, "Failed to get stats", err)
		return
	}

	s.sendSuccessResponse(w, "Statistics retrieved successfully", stats)
}

// fetchRewards 觸發獎勵抓取
func (s *APIServer) fetchRewards(w http.ResponseWriter, r *http.Request) {
	// 建立獎勵抓取器
	fetcher, err := NewRewardFetcher()
	if err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, "Failed to initialize fetcher", err)
		return
	}

	// 抓取獎勵
	rewards, err := fetcher.FetchRewards()
	if err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, "Failed to fetch rewards", err)
		return
	}

	// 儲存到數據庫
	err = s.db.InsertRewards(rewards)
	if err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, "Failed to save rewards", err)
		return
	}

	response := map[string]interface{}{
		"fetched_count": len(rewards),
		"rewards":       rewards,
	}

	s.sendSuccessResponse(w, fmt.Sprintf("Successfully fetched and saved %d rewards", len(rewards)), response)
}

// clearRewards 清空所有獎勵
func (s *APIServer) clearRewards(w http.ResponseWriter, r *http.Request) {
	err := s.db.ClearRewards()
	if err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, "Failed to clear rewards", err)
		return
	}

	s.sendSuccessResponse(w, "All rewards cleared successfully", nil)
}

// healthCheck 健康檢查
func (s *APIServer) healthCheck(w http.ResponseWriter, r *http.Request) {
	health := map[string]string{
		"status":    "healthy",
		"service":   "Flow Rewards API",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	}

	s.sendSuccessResponse(w, "Service is healthy", health)
}

// sendSuccessResponse 發送成功響應
func (s *APIServer) sendSuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// sendErrorResponse 發送錯誤響應
func (s *APIServer) sendErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	response := APIResponse{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Start 啟動 API 服務器
func (s *APIServer) Start() error {
	// 設置 CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // 在生產環境中應該限制具體的域名
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(s.router)

	fmt.Printf("🚀 API Server starting on port %s\n", API_PORT)
	fmt.Printf("📋 API endpoints:\n")
	fmt.Printf("  GET    /api/rewards              - Get all rewards\n")
	fmt.Printf("  GET    /api/rewards/type/{type}  - Get rewards by type\n")
	fmt.Printf("  GET    /api/rewards/delegator/{id} - Get rewards by delegator\n")
	fmt.Printf("  GET    /api/rewards/stats        - Get reward statistics\n")
	fmt.Printf("  POST   /api/fetch               - Fetch new rewards\n")
	fmt.Printf("  DELETE /api/clear               - Clear all rewards\n")
	fmt.Printf("  GET    /api/health              - Health check\n")
	fmt.Printf("\n✅ Server ready at http://localhost:%s\n", API_PORT)

	return http.ListenAndServe(":"+API_PORT, handler)
}

// Close 關閉資源
func (s *APIServer) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func main() {
	server, err := NewAPIServer()
	if err != nil {
		log.Fatalf("❌ Failed to create API server: %v", err)
	}
	defer server.Close()

	if err := server.Start(); err != nil {
		log.Fatalf("❌ Failed to start API server: %v", err)
	}
}
