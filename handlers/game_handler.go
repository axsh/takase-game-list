package handlers

import (
	"net/http"
	"strconv"

	"takase-game-list/models"
	"takase-game-list/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateGame ゲーム登録ハンドラー
func CreateGame(c *gin.Context, db *gorm.DB) {
	var game models.Game

	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateGame(game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create game"})
		return
	}

	c.JSON(http.StatusCreated, game)
}

// GetGames ゲーム一覧取得ハンドラー
func GetGames(c *gin.Context, db *gorm.DB) {
	var games []models.Game
	query := db.Model(&models.Game{})

	// フィルタリング
	if platform := c.Query("platform"); platform != "" {
		query = query.Where("platform = ?", platform)
	}
	if publisher := c.Query("publisher"); publisher != "" {
		query = query.Where("publisher = ?", publisher)
	}
	if genre := c.Query("genre"); genre != "" {
		query = query.Where("genre = ?", genre)
	}
	if series := c.Query("series"); series != "" {
		query = query.Where("series = ?", series)
	}
	if minYear := c.Query("min_year"); minYear != "" {
		query = query.Where("release_year >= ?", minYear)
	}
	if maxYear := c.Query("max_year"); maxYear != "" {
		query = query.Where("release_year <= ?", maxYear)
	}

	// ソート
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

	if order == "asc" {
		query = query.Order(sort + " ASC")
	} else {
		query = query.Order(sort + " DESC")
	}

	if err := query.Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get games"})
		return
	}

	c.JSON(http.StatusOK, games)
}

// SearchGames ゲーム検索ハンドラー
func SearchGames(c *gin.Context, db *gorm.DB) {
	var games []models.Game
	query := db.Model(&models.Game{})

	// 検索キーワード
	keyword := c.Query("q")
	if keyword != "" {
		// タイトル部分一致検索（大文字・小文字を区別しない）
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+keyword+"%")
	}

	// ソート（デフォルト: 登録日時降順）
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

	if order == "asc" {
		query = query.Order(sort + " ASC")
	} else {
		query = query.Order(sort + " DESC")
	}

	if err := query.Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search games"})
		return
	}

	c.JSON(http.StatusOK, games)
}

// UpdateGame ゲーム更新ハンドラー
func UpdateGame(c *gin.Context, db *gorm.DB) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}

	var game models.Game
	if err := db.First(&game, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get game"})
		return
	}

	// 部分更新: 指定されたフィールドのみ更新
	// JSONで送信されたフィールドのみを更新するため、mapで受け取る
	var updateMap map[string]interface{}
	if err := c.ShouldBindJSON(&updateMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 各フィールドが存在する場合のみ更新
	if title, ok := updateMap["title"].(string); ok {
		if title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title cannot be empty"})
			return
		}
		game.Title = title
	}
	if releaseYear, ok := updateMap["release_year"].(float64); ok {
		game.ReleaseYear = int(releaseYear)
	}
	if publisher, ok := updateMap["publisher"].(string); ok {
		if publisher == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "publisher cannot be empty"})
			return
		}
		game.Publisher = publisher
	}
	if platform, ok := updateMap["platform"].(string); ok {
		if platform == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "platform cannot be empty"})
			return
		}
		game.Platform = platform
	}
	if series, ok := updateMap["series"].(string); ok {
		game.Series = &series
	} else if _, ok := updateMap["series"]; ok && updateMap["series"] == nil {
		game.Series = nil
	}
	if genre, ok := updateMap["genre"].(string); ok {
		game.Genre = &genre
	} else if _, ok := updateMap["genre"]; ok && updateMap["genre"] == nil {
		game.Genre = nil
	}
	if price, ok := updateMap["price"].(float64); ok {
		game.Price = int(price)
	}

	// バリデーション
	if err := utils.ValidateGame(game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update game"})
		return
	}

	c.JSON(http.StatusOK, game)
}

// DeleteGame ゲーム削除ハンドラー
func DeleteGame(c *gin.Context, db *gorm.DB) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}

	var game models.Game
	if err := db.First(&game, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get game"})
		return
	}

	if err := db.Delete(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete game"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetStatistics 統計情報取得ハンドラー
func GetStatistics(c *gin.Context, db *gorm.DB) {
	var totalCount int64
	if err := db.Model(&models.Game{}).Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get total count"})
		return
	}

	// プラットフォーム別の集計
	type PlatformCount struct {
		Platform string
		Count    int64
	}
	var platformCounts []PlatformCount
	if err := db.Model(&models.Game{}).
		Select("platform, COUNT(*) as count").
		Group("platform").
		Scan(&platformCounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get platform counts"})
		return
	}
	platformCountsMap := make(map[string]int)
	for _, pc := range platformCounts {
		platformCountsMap[pc.Platform] = int(pc.Count)
	}

	// 発売会社別の集計
	type PublisherCount struct {
		Publisher string
		Count     int64
	}
	var publisherCounts []PublisherCount
	if err := db.Model(&models.Game{}).
		Select("publisher, COUNT(*) as count").
		Group("publisher").
		Scan(&publisherCounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get publisher counts"})
		return
	}
	publisherCountsMap := make(map[string]int)
	for _, pc := range publisherCounts {
		publisherCountsMap[pc.Publisher] = int(pc.Count)
	}

	// ジャンル別の集計（null除外）
	type GenreCount struct {
		Genre string
		Count int64
	}
	var genreCounts []GenreCount
	if err := db.Model(&models.Game{}).
		Select("genre, COUNT(*) as count").
		Where("genre IS NOT NULL AND genre != ''").
		Group("genre").
		Scan(&genreCounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get genre counts"})
		return
	}
	genreCountsMap := make(map[string]int)
	for _, gc := range genreCounts {
		genreCountsMap[gc.Genre] = int(gc.Count)
	}

	// シリーズ別の集計（null除外）
	type SeriesCount struct {
		Series string
		Count  int64
	}
	var seriesCounts []SeriesCount
	if err := db.Model(&models.Game{}).
		Select("series, COUNT(*) as count").
		Where("series IS NOT NULL AND series != ''").
		Group("series").
		Scan(&seriesCounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get series counts"})
		return
	}
	seriesCountsMap := make(map[string]int)
	for _, sc := range seriesCounts {
		seriesCountsMap[sc.Series] = int(sc.Count)
	}

	// 発売年別の集計
	type YearCount struct {
		Year  int
		Count int64
	}
	var yearCounts []YearCount
	if err := db.Model(&models.Game{}).
		Select("release_year as year, COUNT(*) as count").
		Group("release_year").
		Scan(&yearCounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get year counts"})
		return
	}
	yearCountsMap := make(map[int]int)
	for _, yc := range yearCounts {
		yearCountsMap[yc.Year] = int(yc.Count)
	}

	// 価格の合計と平均（価格0の除外）
	var priceStats struct {
		Total int64
		Count int64
	}
	if err := db.Model(&models.Game{}).
		Select("COALESCE(SUM(price), 0) as total, COUNT(*) as count").
		Where("price > 0").
		Scan(&priceStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get price statistics"})
		return
	}

	var averagePrice float64
	if priceStats.Count > 0 {
		averagePrice = float64(priceStats.Total) / float64(priceStats.Count)
	}

	statistics := gin.H{
		"total_count":      int(totalCount),
		"platform_counts":  platformCountsMap,
		"publisher_counts": publisherCountsMap,
		"genre_counts":     genreCountsMap,
		"series_counts":    seriesCountsMap,
		"year_counts":      yearCountsMap,
		"total_price":      int(priceStats.Total),
		"average_price":    averagePrice,
	}

	c.JSON(http.StatusOK, statistics)
}
