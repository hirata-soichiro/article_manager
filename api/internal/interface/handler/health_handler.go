package handler

import (
	"encoding/json"
	"net/http"
)

// HealthResponse はヘルスチェックのレスポンス構造体
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthCheckHandler はヘルスチェックエンドポイントのハンドラー
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status: "healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
