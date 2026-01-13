package models

import (
	"testing"
)

// TestPublisher_Unit Publisher構造体の基本動作確認
func TestPublisher_Unit(t *testing.T) {
	publisher := Publisher{
		Name: "Test Publisher",
	}

	if publisher.Name != "Test Publisher" {
		t.Errorf("expected Name 'Test Publisher', got '%s'", publisher.Name)
	}

	// IDは初期値が0であることを確認
	if publisher.ID != 0 {
		t.Errorf("expected ID 0 (unset), got %d", publisher.ID)
	}
}
