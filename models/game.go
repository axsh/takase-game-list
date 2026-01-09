package models

import (
	"time"

	"gorm.io/gorm"
)

// Game ゲームソフト情報を表すモデル（正規化後）
type Game struct {
	// システム項目
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 必須項目
	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	ReleaseYear int    `gorm:"not null" json:"release_year"`
	PublisherID uint   `gorm:"not null;index" json:"publisher_id"`

	// 任意項目
	SeriesID *uint `gorm:"index" json:"series_id,omitempty"`
	Price    int   `gorm:"default:0" json:"price"`

	// 多対多リレーション
	Platforms []Platform `gorm:"many2many:game_platforms;" json:"platforms,omitempty"`
	Genres    []Genre    `gorm:"many2many:game_genres;" json:"genres,omitempty"`

	// リレーション
	Publisher Publisher `gorm:"foreignKey:PublisherID" json:"publisher,omitempty"`
	Series    *Series   `gorm:"foreignKey:SeriesID" json:"series,omitempty"`
}
