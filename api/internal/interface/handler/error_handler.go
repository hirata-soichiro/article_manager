package handler

import (
	"encoding/json"
	"net/http"

	domainerrors "article-manager/internal/domain/errors"
	"article-manager/internal/domain/service"
	"article-manager/internal/infrastructure/logger"

	"go.uber.org/zap"
)

// エラーレスポンスの構造体
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Code    string                 `json:"code,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ドメインエラーを処理し、適切なHTTPレスポンスを返す
func HandleError(w http.ResponseWriter, err error, operation string) {
	// エラーの種類に応じてステータスコードとメッセージを決定
	statusCode, errorResponse := mapErrorToResponse(err)

	// ログ出力
	logError(err, operation, statusCode)

	// レスポンスを返す
	respondJSON(w, statusCode, errorResponse)
}

// ドメインエラーをHTTPステータスコードとエラーレスポンスにマッピング
func mapErrorToResponse(err error) (int, ErrorResponse) {
	// DomainErrorの場合
	if domainErr, ok := err.(*domainerrors.DomainError); ok {
		statusCode := getStatusCodeFromErrorCode(domainErr.Code)
		return statusCode, ErrorResponse{
			Error:   domainErr.Message,
			Code:    string(domainErr.Code),
			Details: domainErr.Context,
		}
	}

	// AIGeneratorErrorの場合
	if aiErr, ok := err.(*service.AIGeneratorError); ok {
		statusCode := getStatusCodeFromAIError(aiErr)
		return statusCode, ErrorResponse{
			Error: aiErr.Message,
			Code:  string(aiErr.Code),
		}
	}

	// その他のエラー(内部エラーとして扱う)
	return http.StatusInternalServerError, ErrorResponse{
		Error: "internal server error",
		Code:  string(domainerrors.ErrCodeInternal),
	}
}

// エラーコードをHTTPステータスコードにマッピング
func getStatusCodeFromErrorCode(code domainerrors.ErrorCode) int {
	switch code {
	case domainerrors.ErrCodeNotFound:
		return http.StatusNotFound
	case domainerrors.ErrCodeAlreadyExists:
		return http.StatusConflict
	case domainerrors.ErrCodeValidation, domainerrors.ErrCodeInvalidArgument:
		return http.StatusBadRequest
	case domainerrors.ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case domainerrors.ErrCodeForbidden:
		return http.StatusForbidden
	case domainerrors.ErrCodeConflict:
		return http.StatusConflict
	case domainerrors.ErrCodeTimeout:
		return http.StatusGatewayTimeout
	case domainerrors.ErrCodeExternalService:
		return http.StatusBadGateway
	case domainerrors.ErrCodeDatabase, domainerrors.ErrCodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// AIエラーをHTTPステータスコードにマッピング
func getStatusCodeFromAIError(aiErr *service.AIGeneratorError) int {
	switch aiErr.Code {
	case service.ErrCodeAPILimit:
		return http.StatusTooManyRequests
	case service.ErrCodeTimeout:
		return http.StatusGatewayTimeout
	case service.ErrCodeInvalidResponse, service.ErrCodeNetworkError:
		return http.StatusBadGateway
	case service.ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case service.ErrCodeContentBlocked:
		return http.StatusForbidden
	case service.ErrCodeInvalidURL:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// 適切なレベルとコンテキストでエラーをログ出力
func logError(err error, operation string, statusCode int) {
	fields := []zap.Field{
		zap.String("operation", operation),
		zap.Int("status_code", statusCode),
	}

	// DomainErrorの場合、コンテキスト情報を追加
	if domainErr, ok := err.(*domainerrors.DomainError); ok {
		fields = append(
			fields,
			zap.String("error_code", string(domainErr.Code)),
			zap.String("message", domainErr.Message),
		)
		if domainErr.Detail != "" {
			fields = append(fields, zap.String("detail", domainErr.Detail))
		}
		if len(domainErr.Context) > 0 {
			fields = append(fields, zap.Any("context", domainErr.Context))
		}
		if domainErr.OrigErr != nil {
			fields = append(fields, zap.Error(domainErr.OrigErr))
		}
	} else {
		fields = append(fields, zap.Error(err))
	}

	// ステータスコードに応じてログレベルを変更
	switch {
	case statusCode >= 500:
		logger.Error("Request failed with server error", fields...)
	case statusCode >= 400:
		logger.Warn("Request failed with client error", fields...)
	default:
		logger.Info("Request completed", fields...)
	}
}

// JSONレスポンスを送信
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Failed to encode JSON response", zap.Error(err), zap.Int("status_code", statusCode))
	}
}

// 成功したJSONレスポンスを送信
func RespondSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	respondJSON(w, statusCode, data)
}
