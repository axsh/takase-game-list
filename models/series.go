package models

import (
	"time"

	"gorm.io/gorm"
)

// Series シリーズ情報を表すモデル
type Series struct {
	// システム項目
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 必須項目
	Name string `gorm:"type:varchar(255);not null;unique" json:"name"`

	// リレーション
	Games []Game `gorm:"foreignKey:SeriesID" json:"games,omitempty"`
}
