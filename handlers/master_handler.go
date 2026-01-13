package handlers

import (
	"net/http"

	"takase-game-list/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreatePublisherRequest Publisher作成リクエスト構造体
type CreatePublisherRequest struct {
	Name string `json:"name" binding:"required"`
}

// GetPublishers 全Publisherの一覧取得
func GetPublishers(c *gin.Context, db *gorm.DB) {
	var publishers []models.Publisher
	if err := db.Find(&publishers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get publishers"})
		return
	}

	c.JSON(http.StatusOK, publishers)
}

// CreatePublisher Publisherの作成
func CreatePublisher(c *gin.Context, db *gorm.DB) {
	var req CreatePublisherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// バリデーション
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if len(req.Name) > 255 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name must be 255 characters or less"})
		return
	}

	// 重複チェック
	var existingPublisher models.Publisher
	if err := db.Where("name = ?", req.Name).First(&existingPublisher).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "publisher with this name already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check publisher existence"})
		return
	}

	// Publisherの作成
	publisher := models.Publisher{Name: req.Name}
	if err := db.Create(&publisher).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create publisher"})
		return
	}

	c.JSON(http.StatusCreated, publisher)
}

// CreatePlatformRequest Platform作成リクエスト構造体
type CreatePlatformRequest struct {
	Name string `json:"name" binding:"required"`
}

// GetPlatforms 全Platformの一覧取得
func GetPlatforms(c *gin.Context, db *gorm.DB) {
	var platforms []models.Platform
	if err := db.Find(&platforms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get platforms"})
		return
	}

	c.JSON(http.StatusOK, platforms)
}

// CreatePlatform Platformの作成
func CreatePlatform(c *gin.Context, db *gorm.DB) {
	var req CreatePlatformRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// バリデーション
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if len(req.Name) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name must be 100 characters or less"})
		return
	}

	// 重複チェック
	var existingPlatform models.Platform
	if err := db.Where("name = ?", req.Name).First(&existingPlatform).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "platform with this name already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check platform existence"})
		return
	}

	// Platformの作成
	platform := models.Platform{Name: req.Name}
	if err := db.Create(&platform).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create platform"})
		return
	}

	c.JSON(http.StatusCreated, platform)
}

// CreateSeriesRequest Series作成リクエスト構造体
type CreateSeriesRequest struct {
	Name string `json:"name" binding:"required"`
}

// GetSeries 全Seriesの一覧取得
func GetSeries(c *gin.Context, db *gorm.DB) {
	var series []models.Series
	if err := db.Find(&series).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get series"})
		return
	}

	c.JSON(http.StatusOK, series)
}

// CreateSeries Seriesの作成
func CreateSeries(c *gin.Context, db *gorm.DB) {
	var req CreateSeriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// バリデーション
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if len(req.Name) > 255 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name must be 255 characters or less"})
		return
	}

	// 重複チェック
	var existingSeries models.Series
	if err := db.Where("name = ?", req.Name).First(&existingSeries).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "series with this name already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check series existence"})
		return
	}

	// Seriesの作成
	series := models.Series{Name: req.Name}
	if err := db.Create(&series).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create series"})
		return
	}

	c.JSON(http.StatusCreated, series)
}

// CreateGenreRequest Genre作成リクエスト構造体
type CreateGenreRequest struct {
	Name string `json:"name" binding:"required"`
}

// GetGenres 全Genreの一覧取得
func GetGenres(c *gin.Context, db *gorm.DB) {
	var genres []models.Genre
	if err := db.Find(&genres).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get genres"})
		return
	}

	c.JSON(http.StatusOK, genres)
}

// CreateGenre Genreの作成
func CreateGenre(c *gin.Context, db *gorm.DB) {
	var req CreateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// バリデーション
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if len(req.Name) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name must be 100 characters or less"})
		return
	}

	// 重複チェック
	var existingGenre models.Genre
	if err := db.Where("name = ?", req.Name).First(&existingGenre).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "genre with this name already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check genre existence"})
		return
	}

	// Genreの作成
	genre := models.Genre{Name: req.Name}
	if err := db.Create(&genre).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create genre"})
		return
	}

	c.JSON(http.StatusCreated, genre)
}
