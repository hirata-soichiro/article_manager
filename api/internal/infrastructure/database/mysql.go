package database

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// 接続プールを設定
func NewMySQL(host, port, user, password, dbname string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&Asia%%2FTokyo", user, password, host, port, dbname)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %s", err)
	}

	db.SetMaxOpenConns(25)                 // 最大オープン接続数
	db.SetMaxIdleConns(10)                 // アイドル接続の最大数
	db.SetConnMaxLifetime(5 * time.Minute) // 接続の最大生存時間
	db.SetConnMaxIdleTime(5 * time.Minute) // アイドル接続の最大保持時間

	return db, nil
}

// データベース接続を安全に閉じる
func CloseMySQLX(db *sqlx.DB) error {
	if db == nil {
		return nil
	}
	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close MySQL connection: %w", err)
	}

	return nil
}

// データベース接続の健全性を確認する
func HealthCheckX(db *sqlx.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}
