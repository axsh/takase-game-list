package db

import (
	"os"
	"path/filepath"
	"testing"

	"takase-game-list/models"
)

func TestInitDB_Unit(t *testing.T) {
	// 一時的なデータベースファイルを作成
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// データベース接続の初期化
	db, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}

	if db == nil {
		t.Fatal("database instance is nil")
	}

	// 接続の確認
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get database connection: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("database ping failed: %v", err)
	}

	// 接続を閉じる
	sqlDB.Close()

	// データベースファイルが作成されているか確認
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("database file was not created at %s", dbPath)
	}
}

func TestMigrate_Unit(t *testing.T) {
	// 一時的なデータベースファイルを作成
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_migrate.db")

	// データベース接続の初期化
	db, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	// マイグレーション実行
	if err := Migrate(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Gameテーブルが存在するか確認
	if !db.Migrator().HasTable(&models.Game{}) {
		t.Fatal("games table does not exist after migration")
	}
}

func TestGetDB_Unit(t *testing.T) {
	// 初期化されていない状態でGetDBを呼び出すとエラーになることを確認
	// ただし、log.Fatalが呼ばれるため、実際のテストでは初期化後に呼び出す

	// 一時的なデータベースファイルを作成
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_getdb.db")

	// データベース接続の初期化
	_, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}

	// GetDBが正常に動作することを確認
	dbInstance := GetDB()
	if dbInstance == nil {
		t.Fatal("GetDB returned nil")
	}

	// 接続の確認
	sqlDB, err := dbInstance.DB()
	if err != nil {
		t.Fatalf("failed to get database connection: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("database ping failed: %v", err)
	}

	sqlDB.Close()
}
