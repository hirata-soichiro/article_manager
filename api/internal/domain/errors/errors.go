package errors

import (
	"fmt"
)

type ErrorCode string

const (
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists   ErrorCode = "ALREADY_EXISTS"
	ErrCodeValidation      ErrorCode = "VALIDATION"
	ErrCodeInvalidArgument ErrorCode = "INVALID_ARGUMENT"
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrCodeInternal        ErrorCode = "INTERNAL"
	ErrCodeDatabase        ErrorCode = "DATABASE"
	ErrCodeExternalService ErrorCode = "EXTERNAL_SERVICE"
	ErrCodeTimeout         ErrorCode = "TIMEOUT"
	ErrCodeConflict        ErrorCode = "CONFLICT"
)

// ドメイン固有のエラー型
type DomainError struct {
	Code    ErrorCode              // エラーコード
	Message string                 // ユーザー向けメッセージ
	Detail  string                 // 詳細情報
	OrigErr error                  // 元のエラー
	Context map[string]interface{} // コンテキスト情報
}

// インターフェース実装
func (e *DomainError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// エラーにコンテキスト情報を追加
func (e *DomainError) AddContext(key string, value interface{}) *DomainError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// 新しいドメインエラーを作成
func NewDomainError(code ErrorCode, message string, detail string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Detail:  detail,
		Context: make(map[string]interface{}),
	}
}

// 既存エラーをDomainErrorでラップ
func WrapError(code ErrorCode, message string, err error) *DomainError {
	detail := ""
	if err != nil {
		detail = err.Error()
	}
	return &DomainError{
		Code:    code,
		Message: message,
		Detail:  detail,
		OrigErr: err,
		Context: make(map[string]interface{}),
	}
}

// --- コンストラクタ関数 ---

// NotFoundエラーを作成
func NotFoundError(resource string, identifier interface{}) *DomainError {
	return NewDomainError(
		ErrCodeNotFound,
		fmt.Sprintf("%s not found", resource),
		fmt.Sprintf("identifier: %v", identifier),
	).AddContext("resource", resource).AddContext("identifier", identifier)
}

// すでに存在するエラーを作成
func AlreadyExistsError(resource string, identifier interface{}) *DomainError {
	return NewDomainError(
		ErrCodeAlreadyExists,
		fmt.Sprintf("%s already exists", resource),
		fmt.Sprintf("identifier: %v", identifier),
	).AddContext("resource", resource).AddContext("identifier", identifier)
}

// バリデーションエラーを作成
func ValidationError(field string, reason string) *DomainError {
	return NewDomainError(
		ErrCodeValidation,
		fmt.Sprintf("validation failed for %s", field),
		reason,
	).AddContext("field", field).AddContext("reason", reason)
}

// 不正な引数エラーを作成
func InvalidArgumentError(argument string, reason string) *DomainError {
	return NewDomainError(
		ErrCodeInvalidArgument,
		fmt.Sprintf("invalid argument %s", argument),
		reason,
	).AddContext("argument", argument).AddContext("reason", reason)
}

// 内部エラーventbridge作成
func InternalError(message string, err error) *DomainError {
	return WrapError(ErrCodeInternal, message, err)
}

// データベースエラーを作成
func DatabaseError(operation string, err error) *DomainError {
	return WrapError(
		ErrCodeDatabase,
		fmt.Sprintf("database operation failed: %s", operation),
		err,
	).AddContext("operation", operation)
}

// 外部サービスエラーを作成
func ExternalServiceError(service string, err error) *DomainError {
	return WrapError(
		ErrCodeExternalService,
		fmt.Sprintf("external service error: %s", service),
		err,
	).AddContext("service", service)
}

// タイムアウトエラーを作成
func TimeoutError(operation string) *DomainError {
	return NewDomainError(
		ErrCodeTimeout,
		"operation timed out",
		operation,
	).AddContext("operation", operation)
}

// 競合エラーを作成
func ConflictError(resource string, reason string) *DomainError {
	return NewDomainError(
		ErrCodeConflict,
		fmt.Sprintf("%s conflict", resource),
		reason,
	).AddContext("resource", resource).AddContext("reason", reason)
}

// --- ヘルパー関数 ---
// エラーがDomainErrorかどうかをチェック
func IsDomainError(err error) bool {
	_, ok := err.(*DomainError)
	return ok
}

// エラーからエラーコードを取得
func GetErrorCode(err error) ErrorCode {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Code
	}
	return ErrCodeInternal
}

// エラーがNotFoundエラーかどうかをチェック
func IsNotFoundError(err error) bool {
	return GetErrorCode(err) == ErrCodeNotFound
}

// エラーがバリデーションエラーかどうかをチェック
func IsValidationError(err error) bool {
	code := GetErrorCode(err)
	return code == ErrCodeValidation || code == ErrCodeInvalidArgument
}

// エラーが既存エラーかどうかをチェック
func IsAlreadyExistsError(err error) bool {
	return GetErrorCode(err) == ErrCodeAlreadyExists
}
