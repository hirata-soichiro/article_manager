package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *zap.Logger        // アプリケーション全体で使用するロガー
	globalSugar  *zap.SugaredLogger // アプリケーション全体で使用するサガードロガー(より使いやすい)
)

// コンテキストにロガーを格納するためのキー
type contextKey string

const (
	loggerKey contextKey = "logger"
)

// グローバルロガーを初期化
func InitLogger(env string) error {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// ログレベルの設定
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		var level zapcore.Level
		if err := level.UnmarshalText([]byte(logLevel)); err == nil {
			config.Level.SetLevel(level)
		}
	}

	logger, err := config.Build(
		zap.AddCallerSkip(1), // スタックトレースを1つ上のレベルに設定
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return err
	}

	globalLogger = logger
	globalSugar = logger.Sugar()

	return nil
}

// グローバルロガーを取得
func Get() *zap.Logger {
	if globalLogger == nil {
		// 初期化されてない場合はデフォルトロガーを利用
		_ = InitLogger("development")
	}
	return globalLogger
}

// グローバルサガードロガーを取得
func Sugar() *zap.SugaredLogger {
	if globalSugar == nil {
		// 初期化されてない場合はデフォルトロガーを利用
		_ = InitLogger("development")
	}
	return globalSugar
}

// コンテキストからロガーを取得
func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return logger
	}
	return Get()
}

// コンテキストにロガーを追加
func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// 追加フィールドを持つロガーを返す
func WithFields(fields ...zap.Field) *zap.Logger {
	return Get().With(fields...)
}

// バッファされたログエントリをフラッシュ
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// --- 便利なヘルパー関数 ---
// デバッグメッセージをログ出力
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// 情報メッセージをログ出力
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// 警告メッセージをログ出力
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// エラーメッセージをログ出力
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// 致命的なメッセージをログ出力して終了
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

// ログフィールドにエラーを追加
func WithError(err error) zap.Field {
	return zap.Error(err)
}

// ログフィールドにリクエストIDを追加
func WithRequestID(requestID string) zap.Field {
	return zap.String("request_id", requestID)
}

// ログフィールドにユーザーIDを追加
func WithUserID(userID int64) zap.Field {
	return zap.Int64("user_id", userID)
}

// ログフィールドに期間を追加
func WithDuration(key string, duration interface{}) zap.Field {
	return zap.Any(key, duration)
}
