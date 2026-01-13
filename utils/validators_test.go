package utils

import (
	"testing"

	"takase-game-list/models"
)

func TestValidateGame_Unit(t *testing.T) {
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
				Publisher:   "Test Publisher",
				Platform:    "PC",
			},
			wantErr: false,
		},
		{
			name: "タイトルが空",
			game: models.Game{
				Title:       "",
				ReleaseYear: 2024,
				Publisher:   "Test Publisher",
				Platform:    "PC",
			},
			wantErr: true,
		},
		{
			name: "発売年が範囲外（1900未満）",
			game: models.Game{
				Title:       "Test Game",
				ReleaseYear: 1899,
				Publisher:   "Test Publisher",
				Platform:    "PC",
			},
			wantErr: true,
		},
		{
			name: "発売年が範囲外（9999超過）",
			game: models.Game{
				Title:       "Test Game",
				ReleaseYear: 10000,
				Publisher:   "Test Publisher",
				Platform:    "PC",
			},
			wantErr: true,
		},
		{
			name: "発売会社が空",
			game: models.Game{
				Title:       "Test Game",
				ReleaseYear: 2024,
				Publisher:   "",
				Platform:    "PC",
			},
			wantErr: true,
		},
		{
			name: "プラットフォームが空",
			game: models.Game{
				Title:       "Test Game",
				ReleaseYear: 2024,
				Publisher:   "Test Publisher",
				Platform:    "",
			},
			wantErr: true,
		},
		{
			name: "価格が負の値",
			game: models.Game{
				Title:       "Test Game",
				ReleaseYear: 2024,
				Publisher:   "Test Publisher",
				Platform:    "PC",
				Price:       -1,
			},
			wantErr: true,
		},
		{
			name: "価格が0（未登録）",
			game: models.Game{
				Title:       "Test Game",
				ReleaseYear: 2024,
				Publisher:   "Test Publisher",
				Platform:    "PC",
				Price:       0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGame(tt.game)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGame() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
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
