package models

import (
	"testing"
)

// TestPlatform_Unit Platform構造体の基本動作確認
func TestPlatform_Unit(t *testing.T) {
	platform := Platform{
		Name: "PC",
	}

	if platform.Name != "PC" {
		t.Errorf("expected Name 'PC', got '%s'", platform.Name)
	}

	// IDは初期値が0であることを確認
	if platform.ID != 0 {
		t.Errorf("expected ID 0 (unset), got %d", platform.ID)
	}
}
