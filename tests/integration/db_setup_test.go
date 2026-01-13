package integration

import (
	"os"
	"path/filepath"
	"testing"

	"takase-game-list/db"
	"takase-game-list/models"

	"gorm.io/gorm"
)

// setupTestDB テスト用のデータベースをセットアップする
func setupTestDB(t *testing.T) (*gorm.DB, string) {
	t.Helper()

	// 一時的なデータベースファイルを作成
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// データベース接続の初期化
	database, err := db.InitDB(dbPath)
	if err != nil {
		t.Fatalf("failed to initialize test database: %v", err)
	}

	// マイグレーション実行
	if err := db.Migrate(database); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return database, dbPath
}

// cleanupTestDB テスト用のデータベースをクリーンアップする
func cleanupTestDB(t *testing.T, database *gorm.DB, dbPath string) {
	t.Helper()
	if database != nil {
		sqlDB, err := database.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
	if dbPath != "" {
		os.Remove(dbPath)
	}
}

func TestDatabaseConnection_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	// データベース接続が正常に確立されているか確認
	sqlDB, err := database.DB()
	if err != nil {
		t.Fatalf("failed to get database connection: %v", err)
	}

	// 接続のPingテスト
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("database ping failed: %v", err)
	}

	// 接続を閉じる
	sqlDB.Close()
}

func TestDatabaseMigration_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	// Gameテーブルが存在するか確認
	var game models.Game
	if !database.Migrator().HasTable(&game) {
		t.Fatal("games table does not exist")
	}

	// テーブルのカラム構造を確認
	columnTypes, err := database.Migrator().ColumnTypes(&game)
	if err != nil {
		t.Fatalf("failed to get column types: %v", err)
	}

	// 必須カラムが存在するか確認
	requiredColumns := map[string]bool{
		"id":           false,
		"title":        false,
		"release_year": false,
		"publisher":    false,
		"platform":     false,
		"created_at":   false,
		"updated_at":   false,
	}

	for _, ct := range columnTypes {
		columnName := ct.Name()
		if _, exists := requiredColumns[columnName]; exists {
			requiredColumns[columnName] = true
		}
	}

	for column, exists := range requiredColumns {
		if !exists {
			t.Errorf("required column '%s' does not exist", column)
		}
	}
}

func TestGameModel_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	// 必須項目のみでGameを作成
	game := models.Game{
		Title:       "Test Game",
		ReleaseYear: 2024,
		Publisher:   "Test Publisher",
		Platform:    "PC",
	}

	// データベースに保存
	if err := database.Create(&game).Error; err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	// IDが自動採番されているか確認
	if game.ID == 0 {
		t.Error("game ID was not auto-generated")
	}

	// CreatedAtとUpdatedAtが設定されているか確認
	if game.CreatedAt.IsZero() {
		t.Error("CreatedAt was not set")
	}
	if game.UpdatedAt.IsZero() {
		t.Error("UpdatedAt was not set")
	}

	// データベースから取得して確認
	var retrievedGame models.Game
	if err := database.First(&retrievedGame, game.ID).Error; err != nil {
		t.Fatalf("failed to retrieve game: %v", err)
	}

	if retrievedGame.Title != "Test Game" {
		t.Errorf("expected title 'Test Game', got '%s'", retrievedGame.Title)
	}
	if retrievedGame.ReleaseYear != 2024 {
		t.Errorf("expected release year 2024, got %d", retrievedGame.ReleaseYear)
	}
	if retrievedGame.Publisher != "Test Publisher" {
		t.Errorf("expected publisher 'Test Publisher', got '%s'", retrievedGame.Publisher)
	}
	if retrievedGame.Platform != "PC" {
		t.Errorf("expected platform 'PC', got '%s'", retrievedGame.Platform)
	}
}

func TestGameModelWithOptionalFields_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	series := "Test Series"
	genre := "RPG"
	price := 5000

	// 全項目を指定してGameを作成
	game := models.Game{
		Title:       "Test Game 2",
		ReleaseYear: 2023,
		Publisher:   "Test Publisher 2",
		Platform:    "Nintendo Switch",
		Series:      &series,
		Genre:       &genre,
		Price:       price,
	}

	// データベースに保存
	if err := database.Create(&game).Error; err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	// データベースから取得して確認
	var retrievedGame models.Game
	if err := database.First(&retrievedGame, game.ID).Error; err != nil {
		t.Fatalf("failed to retrieve game: %v", err)
	}

	if retrievedGame.Series == nil || *retrievedGame.Series != series {
		t.Errorf("expected series '%s', got %v", series, retrievedGame.Series)
	}
	if retrievedGame.Genre == nil || *retrievedGame.Genre != genre {
		t.Errorf("expected genre '%s', got %v", genre, retrievedGame.Genre)
	}
	if retrievedGame.Price != price {
		t.Errorf("expected price %d, got %d", price, retrievedGame.Price)
	}
}
