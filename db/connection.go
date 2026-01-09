package db

import (
	"log"

	"takase-game-list/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var dbInstance *gorm.DB

// InitDB データベース接続を初期化する
func InitDB(databasePath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 接続プール設定
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	dbInstance = db
	return db, nil
}

// Migrate データベースのマイグレーションを実行する
func Migrate(db *gorm.DB) error {
	// 新規モデルのマイグレーション
	if err := db.AutoMigrate(&models.Publisher{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Platform{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Series{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Genre{}); err != nil {
		return err
	}

	// Gameモデルのマイグレーション（正規化後）
	// 中間テーブル（game_platforms, game_genres）はGORMが自動生成する
	if err := db.AutoMigrate(&models.Game{}); err != nil {
		return err
	}

	return nil
}

// GetDB データベース接続インスタンスを取得する
func GetDB() *gorm.DB {
	if dbInstance == nil {
		log.Fatal("database not initialized. call InitDB first")
	}
	return dbInstance
}
