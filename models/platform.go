package models

import (
	"time"

	"gorm.io/gorm"
)

// Platform プラットフォーム情報を表すモデル
type Platform struct {
	// システム項目
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 必須項目
	Name string `gorm:"type:varchar(100);not null;unique" json:"name"`

	// リレーション
	Games []Game `gorm:"many2many:game_platforms;" json:"games,omitempty"`
}
