package models

import (
	"testing"
)

// TestGenre_Unit Genre構造体の基本動作確認
func TestGenre_Unit(t *testing.T) {
	genre := Genre{
		Name: "RPG",
	}

	if genre.Name != "RPG" {
		t.Errorf("expected Name 'RPG', got '%s'", genre.Name)
	}

	// IDは初期値が0であることを確認
	if genre.ID != 0 {
		t.Errorf("expected ID 0 (unset), got %d", genre.ID)
	}
}
