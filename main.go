package main

import (
	"log"
	"net/http"

	"takase-game-list/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"db":     dbOK(db),
		})
	})
	return router
}

func dbOK(db *gorm.DB) bool {
	sqlDB, err := db.DB()
	if err != nil {
		return false
	}
	return sqlDB.Ping() == nil
}

func main() {
	// データベース接続の初期化
	database, err := db.InitDB("game-list.db")
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// マイグレーション実行
	if err := db.Migrate(database); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	router := setupRouter(database)
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
