package models

import (
	"testing"
)

// TestSeries_Unit Series構造体の基本動作確認
func TestSeries_Unit(t *testing.T) {
	series := Series{
		Name: "Test Series",
	}

	if series.Name != "Test Series" {
		t.Errorf("expected Name 'Test Series', got '%s'", series.Name)
	}

	// IDは初期値が0であることを確認
	if series.ID != 0 {
		t.Errorf("expected ID 0 (unset), got %d", series.ID)
	}
}
