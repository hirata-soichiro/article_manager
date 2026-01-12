package database

import (
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// マイグレーション管理構造体
type MigrationManager struct {
	db     *sqlx.DB
	logger *log.Logger
}

// マイグレーションマネージャー作成
func NewMigrationManager(db *sqlx.DB, logger *log.Logger) *MigrationManager {
	if logger == nil {
		logger = log.Default()
	}
	return &MigrationManager{
		db:     db,
		logger: logger,
	}
}

// マイグレーション実行
func (m *MigrationManager) RunMigrations() error {
	m.logger.Println("マイグレーションを開始します...")

	migrator, err := m.createMigrator()
	if err != nil {
		return fmt.Errorf("マイグレーターの作成に失敗: %w", err)
	}

	// 最新バージョンを確認
	currentVersion, dirty, err := migrator.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("現在のマイグレーションバージョンの取得に失敗: %w", err)
	}

	// Dirty状態の検出
	if dirty {
		return fmt.Errorf("マイグレーションがdirty状態です(バージョン: %d)。手動で修復が必要です", currentVersion)
	}

	if err == migrate.ErrNilVersion {
		m.logger.Println("マイグレーション履歴がありません。初回実行です。")
	} else {
		m.logger.Printf("現在のマイグレーションバージョン: %d\n", currentVersion)
	}

	// 最新バージョンまでマイグレーションを実行
	if err := migrator.Up(); err != nil {
		if err == migrate.ErrNoChange {
			m.logger.Println("適用するマイグレーションはありません。データベースは最新です。")
			return nil
		}
		return fmt.Errorf("マイグレーション実行に失敗: %w", err)
	}

	// 実行後のバージョンを確認
	newVersion, _, err := migrator.Version()
	if err != nil {
		return fmt.Errorf("新しいバージョンの取得に失敗: %w", err)
	}

	m.logger.Printf("マイグレーションが正常に完了しました。現在のバージョン: %d\n", newVersion)
	return nil
}

// 指定した数のマイグレーションをロールバック
func (m *MigrationManager) MigrateDown(steps int) error {
	m.logger.Printf("マイグレーションのロールバックを開始します(ステップ数: %d) ...\n", steps)

	migrator, err := m.createMigrator()
	if err != nil {
		return fmt.Errorf("マイグレーターの作成に失敗: %w", err)
	}

	if steps <= 0 {
		// 全ロールバック
		if err := migrator.Down(); err != nil {
			if err == migrate.ErrNoChange {
				m.logger.Println("ロールバックするマイグレーションはありません。")
				return nil
			}
			return fmt.Errorf("全ロールバックに失敗: %w", err)
		}
		m.logger.Println("全マイグレーションのロールバックが完了しました。")
	} else {
		// 指定ステップ数のロールバック
		if err := migrator.Steps(-steps); err != nil {
			if err == migrate.ErrNoChange {
				m.logger.Println("ロールバックするマイグレーションはありません。")
				return nil
			}
			return fmt.Errorf("ロールバックに失敗: %w", err)
		}
		m.logger.Printf("%dステップのロールバックが完了しました。\n", steps)
	}

	return nil
}

// 現在のマイグレーションバージョンを取得
func (m *MigrationManager) GetCurrentVersion() (version uint, dirty bool, err error) {
	migrator, err := m.createMigrator()
	if err != nil {
		return 0, false, fmt.Errorf("マイグレーターの作成に失敗： %w", err)
	}

	version, dirty, err = migrator.Version()
	if err == migrate.ErrNilVersion {
		return 0, false, nil
	}

	return version, dirty, err
}

// migrateインスタンスを作成
func (m *MigrationManager) createMigrator() (*migrate.Migrate, error) {
	// embed.FSからマイグレーションファイルを読み込む
	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("マイグレーションソースの作成に失敗: %w", err)
	}

	// MySQLドライバー作成
	stdDB := m.db.DB
	driver, err := mysql.WithInstance(stdDB, &mysql.Config{
		MigrationsTable: "schema_migrations",
		DatabaseName:    "",
	})
	if err != nil {
		return nil, fmt.Errorf("MySQLドライバーの作成に失敗： %w", err)
	}

	// migrateインスタンスを作成
	migrator, err := migrate.NewWithInstance("iofs", sourceDriver, "mysql", driver)
	if err != nil {
		return nil, fmt.Errorf("マイグレーターの初期化に失敗： %w", err)
	}

	return migrator, nil
}
