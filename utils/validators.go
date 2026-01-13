package utils

import (
	"errors"
	"fmt"

	"takase-game-list/models"
)

// ValidateGame ゲーム情報のバリデーション
func ValidateGame(game models.Game) error {
	if game.Title == "" {
		return errors.New("title is required")
	}
	if game.Publisher == "" {
		return errors.New("publisher is required")
	}
	if game.Platform == "" {
		return errors.New("platform is required")
	}
	if !ValidateReleaseYear(game.ReleaseYear) {
		return fmt.Errorf("release_year must be between 1900 and 9999, got %d", game.ReleaseYear)
	}
	if !ValidatePrice(game.Price) {
		return fmt.Errorf("price must be 0 or greater, got %d", game.Price)
	}
	return nil
}

// ValidateReleaseYear 発売年の範囲チェック（1900-9999）
func ValidateReleaseYear(year int) bool {
	return year >= 1900 && year <= 9999
}

// ValidatePrice 価格の非負チェック
func ValidatePrice(price int) bool {
	return price >= 0
}
