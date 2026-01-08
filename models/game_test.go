package models

import (
	"testing"
)

func TestGame_Unit(t *testing.T) {
	// Game構造体の基本的な動作確認
	game := Game{
		Title:       "Test Game",
		ReleaseYear: 2024,
		Publisher:   "Test Publisher",
		Platform:    "PC",
	}

	if game.Title != "Test Game" {
		t.Errorf("expected Title 'Test Game', got '%s'", game.Title)
	}
	if game.ReleaseYear != 2024 {
		t.Errorf("expected ReleaseYear 2024, got %d", game.ReleaseYear)
	}
	if game.Publisher != "Test Publisher" {
		t.Errorf("expected Publisher 'Test Publisher', got '%s'", game.Publisher)
	}
	if game.Platform != "PC" {
		t.Errorf("expected Platform 'PC', got '%s'", game.Platform)
	}
}

func TestGame_WithOptionalFields_Unit(t *testing.T) {
	series := "Test Series"
	genre := "RPG"
	price := 5000

	game := Game{
		Title:       "Test Game 2",
		ReleaseYear: 2023,
		Publisher:   "Test Publisher 2",
		Platform:    "Nintendo Switch",
		Series:      &series,
		Genre:       &genre,
		Price:       price,
	}

	if game.Series == nil || *game.Series != series {
		t.Errorf("expected Series '%s', got %v", series, game.Series)
	}
	if game.Genre == nil || *game.Genre != genre {
		t.Errorf("expected Genre '%s', got %v", genre, game.Genre)
	}
	if game.Price != price {
		t.Errorf("expected Price %d, got %d", price, game.Price)
	}
}

func TestGame_WithNilOptionalFields_Unit(t *testing.T) {
	// 任意項目がnilの場合の確認
	game := Game{
		Title:       "Test Game 3",
		ReleaseYear: 2022,
		Publisher:   "Test Publisher 3",
		Platform:    "PlayStation 5",
		Series:      nil,
		Genre:       nil,
		Price:       0,
	}

	if game.Series != nil {
		t.Error("expected Series to be nil")
	}
	if game.Genre != nil {
		t.Error("expected Genre to be nil")
	}
	if game.Price != 0 {
		t.Errorf("expected Price 0 (unset), got %d", game.Price)
	}
}
