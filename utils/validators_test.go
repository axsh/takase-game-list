package utils

import (
	"testing"

	"takase-game-list/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// setupTestDBForValidators バリデーションテスト用のテストデータベースをセットアップ
func setupTestDBForValidators(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// マイグレーション実行
	if err := db.AutoMigrate(&models.Publisher{}, &models.Platform{}, &models.Series{}, &models.Genre{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	// テストデータの投入
	publisher := models.Publisher{Name: "Test Publisher"}
	db.Create(&publisher)

	platform1 := models.Platform{Name: "PC"}
	platform2 := models.Platform{Name: "Nintendo Switch"}
	db.Create(&platform1)
	db.Create(&platform2)

	series := models.Series{Name: "Test Series"}
	db.Create(&series)

	genre1 := models.Genre{Name: "RPG"}
	genre2 := models.Genre{Name: "Action"}
	db.Create(&genre1)
	db.Create(&genre2)

	return db
}

func TestValidateGame_Unit(t *testing.T) {
	db := setupTestDBForValidators(t)

	tests := []struct {
		name    string
		game    models.Game
		wantErr bool
	}{
		{
			name: "正常なゲーム情報",
			game: models.Game{
				Title:       "Test Game",
				ReleaseYear: 2024,
				PublisherID: 1,
				Platforms: []models.Platform{
					{ID: 1},
				},
			},
			wantErr: false,
		},
		{
			name: "タイトルが空",
			game: models.Game{
				Title:       "",
				ReleaseYear: 2024,
				PublisherID: 1,
			},
			wantErr: true,
		},
		{
			name: "発売年が範囲外（1900未満）",
			game: models.Game{
				Title:       "Test Game",
				ReleaseYear: 1899,
				PublisherID: 1,
			},
			wantErr: true,
		},
		{
			name: "発売年が範囲外（9999超過）",
			game: models.Game{
				Title:       "Test Game",
				ReleaseYear: 10000,
				PublisherID: 1,
			},
			wantErr: true,
		},
		{
			name: "PublisherIDが0（未設定）",
			game: models.Game{
				Title:       "Test Game",
				ReleaseYear: 2024,
				PublisherID: 0,
			},
			wantErr: true,
		},
		{
			name: "PublisherIDが存在しない",
			game: models.Game{
				Title:       "Test Game",
				ReleaseYear: 2024,
				PublisherID: 999,
			},
			wantErr: true,
		},
		{
			name: "PlatformIDsが空配列",
			game: models.Game{
				Title:       "Test Game",
				ReleaseYear: 2024,
				PublisherID: 1,
			},
			wantErr: true,
		},
		{
			name: "PlatformIDsに存在しないIDが含まれる",
			game: models.Game{
				Title:       "Test Game",
				ReleaseYear: 2024,
				PublisherID: 1,
				Platforms: []models.Platform{
					{ID: 999},
				},
			},
			wantErr: true,
		},
		{
			name: "価格が負の値",
			game: models.Game{
				Title:       "Test Game",
				ReleaseYear: 2024,
				PublisherID: 1,
				Platforms: []models.Platform{
					{ID: 1},
				},
				Price: -1,
			},
			wantErr: true,
		},
		{
			name: "価格が0（未登録）",
			game: models.Game{
				Title:       "Test Game",
				ReleaseYear: 2024,
				PublisherID: 1,
				Platforms: []models.Platform{
					{ID: 1},
				},
				Price: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGame(tt.game, db)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGame() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidatePublisherID_Unit ValidatePublisherID関数のテスト
func TestValidatePublisherID_Unit(t *testing.T) {
	db := setupTestDBForValidators(t)

	tests := []struct {
		name      string
		publisherID uint
		wantErr   bool
	}{
		{"存在するPublisherID", 1, false},
		{"存在しないPublisherID", 999, true},
		{"PublisherIDが0", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePublisherID(tt.publisherID, db)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePublisherID(%d) error = %v, wantErr %v", tt.publisherID, err, tt.wantErr)
			}
		})
	}
}

// TestValidateSeriesID_Unit ValidateSeriesID関数のテスト
func TestValidateSeriesID_Unit(t *testing.T) {
	db := setupTestDBForValidators(t)

	tests := []struct {
		name    string
		seriesID *uint
		wantErr bool
	}{
		{"存在するSeriesID", uintPtr(1), false},
		{"存在しないSeriesID", uintPtr(999), true},
		{"SeriesIDがnil（許容）", nil, false},
		{"SeriesIDが0", uintPtr(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSeriesID(tt.seriesID, db)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSeriesID(%v) error = %v, wantErr %v", tt.seriesID, err, tt.wantErr)
			}
		})
	}
}

// TestValidatePlatformIDs_Unit ValidatePlatformIDs関数のテスト
func TestValidatePlatformIDs_Unit(t *testing.T) {
	db := setupTestDBForValidators(t)

	tests := []struct {
		name       string
		platformIDs []uint
		wantErr    bool
	}{
		{"正常なPlatformIDs（1つ）", []uint{1}, false},
		{"正常なPlatformIDs（複数）", []uint{1, 2}, false},
		{"空配列（エラー）", []uint{}, true},
		{"存在しないIDが含まれる", []uint{1, 999}, true},
		{"全て存在しないID", []uint{999, 998}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePlatformIDs(tt.platformIDs, db)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePlatformIDs(%v) error = %v, wantErr %v", tt.platformIDs, err, tt.wantErr)
			}
		})
	}
}

// TestValidateGenreIDs_Unit ValidateGenreIDs関数のテスト
func TestValidateGenreIDs_Unit(t *testing.T) {
	db := setupTestDBForValidators(t)

	tests := []struct {
		name     string
		genreIDs []uint
		wantErr  bool
	}{
		{"正常なGenreIDs（1つ）", []uint{1}, false},
		{"正常なGenreIDs（複数）", []uint{1, 2}, false},
		{"空配列（許容）", []uint{}, false},
		{"存在しないIDが含まれる", []uint{1, 999}, true},
		{"全て存在しないID", []uint{999, 998}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGenreIDs(tt.genreIDs, db)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGenreIDs(%v) error = %v, wantErr %v", tt.genreIDs, err, tt.wantErr)
			}
		})
	}
}

// uintPtr uintのポインタを返すヘルパー関数
func uintPtr(u uint) *uint {
	return &u
}

func TestValidateReleaseYear_Unit(t *testing.T) {
	tests := []struct {
		name  string
		year  int
		valid bool
	}{
		{"正常な年（1900）", 1900, true},
		{"正常な年（2024）", 2024, true},
		{"正常な年（9999）", 9999, true},
		{"範囲外（1899）", 1899, false},
		{"範囲外（10000）", 10000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := ValidateReleaseYear(tt.year)
			if valid != tt.valid {
				t.Errorf("ValidateReleaseYear(%d) = %v, want %v", tt.year, valid, tt.valid)
			}
		})
	}
}

func TestValidatePrice_Unit(t *testing.T) {
	tests := []struct {
		name  string
		price int
		valid bool
	}{
		{"正常な価格（0）", 0, true},
		{"正常な価格（1000）", 1000, true},
		{"正常な価格（10000）", 10000, true},
		{"負の値（-1）", -1, false},
		{"負の値（-100）", -100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := ValidatePrice(tt.price)
			if valid != tt.valid {
				t.Errorf("ValidatePrice(%d) = %v, want %v", tt.price, valid, tt.valid)
			}
		})
	}
}
