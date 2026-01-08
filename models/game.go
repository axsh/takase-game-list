package models

import (
	"time"

	"gorm.io/gorm"
)

// Game ゲームソフト情報を表すモデル
type Game struct {
	// システム項目
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 必須項目
	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	ReleaseYear int    `gorm:"not null" json:"release_year"`
	Publisher   string `gorm:"type:varchar(255);not null" json:"publisher"`
	Platform    string `gorm:"type:varchar(100);not null" json:"platform"`

	// 任意項目
	Series *string `gorm:"type:varchar(255)" json:"series,omitempty"`
	Genre  *string `gorm:"type:varchar(100)" json:"genre,omitempty"`
	Price  int     `gorm:"default:0" json:"price"`
}
