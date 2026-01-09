package main

import (
	"log"
	"net/http"

	"takase-game-list/db"
	"takase-game-list/handlers"
	"takase-game-list/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// ミドルウェアの適用
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())

	// ヘルスチェック
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"db":     dbOK(db),
		})
	})

	// ゲームAPI
	router.POST("/games", func(c *gin.Context) {
		handlers.CreateGame(c, db)
	})
	router.GET("/games", func(c *gin.Context) {
		handlers.GetGames(c, db)
	})
	router.PUT("/games/:id", func(c *gin.Context) {
		handlers.UpdateGame(c, db)
	})
	router.DELETE("/games/:id", func(c *gin.Context) {
		handlers.DeleteGame(c, db)
	})
	router.GET("/games/search", func(c *gin.Context) {
		handlers.SearchGames(c, db)
	})
	router.GET("/games/statistics", func(c *gin.Context) {
		handlers.GetStatistics(c, db)
	})

	// マスタデータAPI
	router.GET("/publishers", func(c *gin.Context) {
		handlers.GetPublishers(c, db)
	})
	router.POST("/publishers", func(c *gin.Context) {
		handlers.CreatePublisher(c, db)
	})
	router.GET("/platforms", func(c *gin.Context) {
		handlers.GetPlatforms(c, db)
	})
	router.POST("/platforms", func(c *gin.Context) {
		handlers.CreatePlatform(c, db)
	})
	router.GET("/series", func(c *gin.Context) {
		handlers.GetSeries(c, db)
	})
	router.POST("/series", func(c *gin.Context) {
		handlers.CreateSeries(c, db)
	})
	router.GET("/genres", func(c *gin.Context) {
		handlers.GetGenres(c, db)
	})
	router.POST("/genres", func(c *gin.Context) {
		handlers.CreateGenre(c, db)
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

	// シードデータの投入（オプション）
	if err := db.SeedMasterData(database); err != nil {
		log.Printf("warning: failed to seed master data: %v", err)
	}

	router := setupRouter(database)
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
