package utils

import (
	"errors"
	"fmt"

	"takase-game-list/models"

	"gorm.io/gorm"
)

// ValidateGame ゲーム情報のバリデーション（正規化後）
func ValidateGame(game models.Game, db *gorm.DB) error {
	if game.Title == "" {
		return errors.New("title is required")
	}
	if !ValidateReleaseYear(game.ReleaseYear) {
		return fmt.Errorf("release_year must be between 1900 and 9999, got %d", game.ReleaseYear)
	}
	if !ValidatePrice(game.Price) {
		return fmt.Errorf("price must be 0 or greater, got %d", game.Price)
	}

	// PublisherIDの存在確認（必須）
	if err := ValidatePublisherID(game.PublisherID, db); err != nil {
		return err
	}

	// SeriesIDの存在確認（任意、nilでない場合のみ）
	if game.SeriesID != nil {
		if err := ValidateSeriesID(game.SeriesID, db); err != nil {
			return err
		}
	}

	// Platformsの存在確認（必須、少なくとも1つ必要）
	platformIDs := make([]uint, len(game.Platforms))
	for i, platform := range game.Platforms {
		platformIDs[i] = platform.ID
	}
	if err := ValidatePlatformIDs(platformIDs, db); err != nil {
		return err
	}

	// Genresの存在確認（任意、空配列も許容）
	genreIDs := make([]uint, len(game.Genres))
	for i, genre := range game.Genres {
		genreIDs[i] = genre.ID
	}
	if err := ValidateGenreIDs(genreIDs, db); err != nil {
		return err
	}

	return nil
}

// ValidatePublisherID PublisherIDの存在確認
func ValidatePublisherID(publisherID uint, db *gorm.DB) error {
	if publisherID == 0 {
		return errors.New("publisher_id is required")
	}

	var publisher models.Publisher
	if err := db.First(&publisher, publisherID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("publisher_id %d does not exist", publisherID)
		}
		return fmt.Errorf("failed to validate publisher_id: %w", err)
	}

	return nil
}

// ValidateSeriesID SeriesIDの存在確認（nil許容）
func ValidateSeriesID(seriesID *uint, db *gorm.DB) error {
	if seriesID == nil {
		return nil // nilは許容
	}

	if *seriesID == 0 {
		return errors.New("series_id cannot be 0")
	}

	var series models.Series
	if err := db.First(&series, *seriesID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("series_id %d does not exist", *seriesID)
		}
		return fmt.Errorf("failed to validate series_id: %w", err)
	}

	return nil
}

// ValidatePlatformIDs PlatformID配列の存在確認（少なくとも1つ必要）
func ValidatePlatformIDs(platformIDs []uint, db *gorm.DB) error {
	if len(platformIDs) == 0 {
		return errors.New("at least one platform_id is required")
	}

	for _, platformID := range platformIDs {
		if platformID == 0 {
			return errors.New("platform_id cannot be 0")
		}

		var platform models.Platform
		if err := db.First(&platform, platformID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("platform_id %d does not exist", platformID)
			}
			return fmt.Errorf("failed to validate platform_id: %w", err)
		}
	}

	return nil
}

// ValidateGenreIDs GenreID配列の存在確認（空配列許容）
func ValidateGenreIDs(genreIDs []uint, db *gorm.DB) error {
	// 空配列は許容
	if len(genreIDs) == 0 {
		return nil
	}

	for _, genreID := range genreIDs {
		if genreID == 0 {
			return errors.New("genre_id cannot be 0")
		}

		var genre models.Genre
		if err := db.First(&genre, genreID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("genre_id %d does not exist", genreID)
			}
			return fmt.Errorf("failed to validate genre_id: %w", err)
		}
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
