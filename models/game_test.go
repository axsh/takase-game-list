package models

import (
	"testing"
)

// TestGame_Unit Game構造体の基本動作確認（正規化後）
func TestGame_Unit(t *testing.T) {
	// Game構造体の基本的な動作確認
	game := Game{
		Title:       "Test Game",
		ReleaseYear: 2024,
		PublisherID: 1,
	}

	if game.Title != "Test Game" {
		t.Errorf("expected Title 'Test Game', got '%s'", game.Title)
	}
	if game.ReleaseYear != 2024 {
		t.Errorf("expected ReleaseYear 2024, got %d", game.ReleaseYear)
	}
	if game.PublisherID != 1 {
		t.Errorf("expected PublisherID 1, got %d", game.PublisherID)
	}

	// IDは初期値が0であることを確認
	if game.ID != 0 {
		t.Errorf("expected ID 0 (unset), got %d", game.ID)
	}
}

// TestGameWithRelations_Unit リレーションフィールドの確認
func TestGameWithRelations_Unit(t *testing.T) {
	publisherID := uint(1)
	seriesID := uint(2)

	game := Game{
		Title:       "Test Game 2",
		ReleaseYear: 2023,
		PublisherID: publisherID,
		SeriesID:    &seriesID,
		Platforms: []Platform{
			{ID: 1, Name: "PC"},
			{ID: 2, Name: "Nintendo Switch"},
		},
		Genres: []Genre{
			{ID: 1, Name: "RPG"},
			{ID: 2, Name: "Action"},
		},
		Price: 5000,
	}

	if game.PublisherID != publisherID {
		t.Errorf("expected PublisherID %d, got %d", publisherID, game.PublisherID)
	}
	if game.SeriesID == nil || *game.SeriesID != seriesID {
		t.Errorf("expected SeriesID %d, got %v", seriesID, game.SeriesID)
	}
	if len(game.Platforms) != 2 {
		t.Errorf("expected 2 platforms, got %d", len(game.Platforms))
	}
	if len(game.Genres) != 2 {
		t.Errorf("expected 2 genres, got %d", len(game.Genres))
	}
	if game.Price != 5000 {
		t.Errorf("expected Price 5000, got %d", game.Price)
	}
}

// TestGameWithOptionalFields_Unit 任意項目の確認
func TestGameWithOptionalFields_Unit(t *testing.T) {
	seriesID := uint(2)
	price := 5000

	game := Game{
		Title:       "Test Game 3",
		ReleaseYear: 2023,
		PublisherID: 1,
		SeriesID:    &seriesID,
		Platforms: []Platform{
			{ID: 1, Name: "PC"},
		},
		Genres: []Genre{
			{ID: 1, Name: "RPG"},
		},
		Price: price,
	}

	if game.SeriesID == nil || *game.SeriesID != seriesID {
		t.Errorf("expected SeriesID %d, got %v", seriesID, game.SeriesID)
	}
	if len(game.Genres) != 1 {
		t.Errorf("expected 1 genre, got %d", len(game.Genres))
	}
	if game.Price != price {
		t.Errorf("expected Price %d, got %d", price, game.Price)
	}
}

// TestGameWithNilOptionalFields_Unit nil許容フィールドの確認
func TestGameWithNilOptionalFields_Unit(t *testing.T) {
	// 任意項目がnilの場合の確認
	game := Game{
		Title:       "Test Game 4",
		ReleaseYear: 2022,
		PublisherID: 1,
		SeriesID:    nil,
		Platforms: []Platform{
			{ID: 1, Name: "PC"},
		},
		Genres: []Genre{},
		Price:  0,
	}

	if game.SeriesID != nil {
		t.Error("expected SeriesID to be nil")
	}
	if len(game.Genres) != 0 {
		t.Errorf("expected 0 genres, got %d", len(game.Genres))
	}
	if game.Price != 0 {
		t.Errorf("expected Price 0 (unset), got %d", game.Price)
	}
}
