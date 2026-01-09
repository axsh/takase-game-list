package db

import (
	"takase-game-list/models"

	"gorm.io/gorm"
)

// SeedMasterData よく使われるマスタデータを投入する
func SeedMasterData(db *gorm.DB) error {
	// Platformデータの投入
	platforms := []models.Platform{
		{Name: "PC"},
		{Name: "PlayStation 5"},
		{Name: "Nintendo Switch"},
		{Name: "Xbox Series X"},
		{Name: "PlayStation 4"},
		{Name: "Xbox One"},
		{Name: "Nintendo 3DS"},
		{Name: "PlayStation Vita"},
	}

	for _, platform := range platforms {
		if err := db.FirstOrCreate(&platform, models.Platform{Name: platform.Name}).Error; err != nil {
			return err
		}
	}

	// Publisherデータの投入
	publishers := []models.Publisher{
		{Name: "Nintendo"},
		{Name: "Sony Interactive Entertainment"},
		{Name: "Microsoft"},
		{Name: "Square Enix"},
		{Name: "Bandai Namco Entertainment"},
		{Name: "Capcom"},
		{Name: "Electronic Arts"},
		{Name: "Ubisoft"},
		{Name: "Activision"},
		{Name: "Sega"},
	}

	for _, publisher := range publishers {
		if err := db.FirstOrCreate(&publisher, models.Publisher{Name: publisher.Name}).Error; err != nil {
			return err
		}
	}

	// Genreデータの投入
	genres := []models.Genre{
		{Name: "アクション"},
		{Name: "RPG"},
		{Name: "アドベンチャー"},
		{Name: "シューティング"},
		{Name: "レーシング"},
		{Name: "スポーツ"},
		{Name: "パズル"},
		{Name: "シミュレーション"},
		{Name: "ストラテジー"},
		{Name: "格闘"},
	}

	for _, genre := range genres {
		if err := db.FirstOrCreate(&genre, models.Genre{Name: genre.Name}).Error; err != nil {
			return err
		}
	}

	return nil
}
