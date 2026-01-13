package handlers

import (
	"net/http"
	"strconv"

	"takase-game-list/models"
	"takase-game-list/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateGameRequest ゲーム作成リクエスト構造体
type CreateGameRequest struct {
	Title       string `json:"title" binding:"required"`
	ReleaseYear int    `json:"release_year" binding:"required"`
	PublisherID uint   `json:"publisher_id" binding:"required"`
	SeriesID    *uint  `json:"series_id,omitempty"`
	PlatformIDs []uint `json:"platform_ids" binding:"required"`
	GenreIDs    []uint `json:"genre_ids,omitempty"`
	Price       int    `json:"price"`
}

// CreateGame ゲーム登録ハンドラー（正規化後）
func CreateGame(c *gin.Context, db *gorm.DB) {
	var req CreateGameRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Platformの取得
	var platforms []models.Platform
	if err := db.Where("id IN ?", req.PlatformIDs).Find(&platforms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get platforms"})
		return
	}
	if len(platforms) != len(req.PlatformIDs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "some platform_ids do not exist"})
		return
	}

	// Genreの取得（空配列の場合はスキップ）
	var genres []models.Genre
	if len(req.GenreIDs) > 0 {
		if err := db.Where("id IN ?", req.GenreIDs).Find(&genres).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get genres"})
			return
		}
		if len(genres) != len(req.GenreIDs) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "some genre_ids do not exist"})
			return
		}
	}

	// Gameの作成
	game := models.Game{
		Title:       req.Title,
		ReleaseYear: req.ReleaseYear,
		PublisherID: req.PublisherID,
		SeriesID:    req.SeriesID,
		Platforms:   platforms,
		Genres:      genres,
		Price:       req.Price,
	}

	// バリデーション
	if err := utils.ValidateGame(game, db); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// データベースに保存
	if err := db.Create(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create game"})
		return
	}

	// Preloadを使用して関連データを取得
	if err := db.Preload("Publisher").Preload("Platforms").Preload("Series").Preload("Genres").
		First(&game, game.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load game relations"})
		return
	}

	c.JSON(http.StatusCreated, game)
}

// GetGames ゲーム一覧取得ハンドラー（正規化後）
func GetGames(c *gin.Context, db *gorm.DB) {
	var games []models.Game
	query := db.Model(&models.Game{})

	// フィルタリング（正規化後の構造に対応）
	hasJoin := false
	if publisherID := c.Query("publisher_id"); publisherID != "" {
		if id, err := strconv.ParseUint(publisherID, 10, 32); err == nil {
			query = query.Where("publisher_id = ?", id)
		}
	}
	if seriesID := c.Query("series_id"); seriesID != "" {
		if id, err := strconv.ParseUint(seriesID, 10, 32); err == nil {
			query = query.Where("series_id = ?", id)
		}
	}
	if platformIDs := c.QueryArray("platform_ids"); len(platformIDs) > 0 {
		var ids []uint
		for _, idStr := range platformIDs {
			if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
				ids = append(ids, uint(id))
			}
		}
		if len(ids) > 0 {
			query = query.Joins("JOIN game_platforms ON games.id = game_platforms.game_id").
				Where("game_platforms.platform_id IN ?", ids)
			hasJoin = true
		}
	}
	if genreIDs := c.QueryArray("genre_ids"); len(genreIDs) > 0 {
		var ids []uint
		for _, idStr := range genreIDs {
			if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
				ids = append(ids, uint(id))
			}
		}
		if len(ids) > 0 {
			query = query.Joins("JOIN game_genres ON games.id = game_genres.game_id").
				Where("game_genres.genre_id IN ?", ids)
			hasJoin = true
		}
	}
	// JOINがある場合はGROUP BYを適用
	if hasJoin {
		query = query.Group("games.id")
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

	// Preloadを使用して関連データを取得
	query = query.Preload("Publisher").Preload("Platforms").Preload("Series").Preload("Genres")

	if err := query.Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get games"})
		return
	}

	c.JSON(http.StatusOK, games)
}

// SearchGames ゲーム検索ハンドラー（正規化後）
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

	// Preloadを使用して関連データを取得
	query = query.Preload("Publisher").Preload("Platforms").Preload("Series").Preload("Genres")

	if err := query.Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search games"})
		return
	}

	c.JSON(http.StatusOK, games)
}

// UpdateGameRequest ゲーム更新リクエスト構造体
type UpdateGameRequest struct {
	Title       *string `json:"title,omitempty"`
	ReleaseYear *int    `json:"release_year,omitempty"`
	PublisherID *uint   `json:"publisher_id,omitempty"`
	SeriesID    *uint   `json:"series_id,omitempty"`
	PlatformIDs []uint  `json:"platform_ids,omitempty"`
	GenreIDs    []uint  `json:"genre_ids,omitempty"`
	Price       *int    `json:"price,omitempty"`
}

// UpdateGame ゲーム更新ハンドラー（正規化後）
func UpdateGame(c *gin.Context, db *gorm.DB) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}

	var game models.Game
	// 既存のGameを取得（PlatformsとGenresもPreload）
	if err := db.Preload("Platforms").Preload("Genres").First(&game, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get game"})
		return
	}

	// まず map でパースして、series_id キーが存在するかどうかを確認
	var updateMap map[string]interface{}
	if err := c.ShouldBindJSON(&updateMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// UpdateGameRequest に手動で変換
	var req UpdateGameRequest
	if title, ok := updateMap["title"].(string); ok {
		req.Title = &title
	}
	if releaseYear, ok := updateMap["release_year"].(float64); ok {
		ry := int(releaseYear)
		req.ReleaseYear = &ry
	}
	if publisherID, ok := updateMap["publisher_id"].(float64); ok {
		pid := uint(publisherID)
		req.PublisherID = &pid
	}
	if seriesID, ok := updateMap["series_id"]; ok {
		if seriesID == nil {
			// null の場合は nil ポインタを設定
			req.SeriesID = nil
		} else if sid, ok := seriesID.(float64); ok {
			sidUint := uint(sid)
			req.SeriesID = &sidUint
		}
	}
	if price, ok := updateMap["price"].(float64); ok {
		p := int(price)
		req.Price = &p
	}
	if platformIDs, ok := updateMap["platform_ids"].([]interface{}); ok {
		req.PlatformIDs = make([]uint, len(platformIDs))
		for i, pid := range platformIDs {
			if pidFloat, ok := pid.(float64); ok {
				req.PlatformIDs[i] = uint(pidFloat)
			}
		}
	}
	if genreIDs, ok := updateMap["genre_ids"].([]interface{}); ok {
		req.GenreIDs = make([]uint, len(genreIDs))
		for i, gid := range genreIDs {
			if gidFloat, ok := gid.(float64); ok {
				req.GenreIDs[i] = uint(gidFloat)
			}
		}
	}

	// 各フィールドが存在する場合のみ更新
	if req.Title != nil {
		if *req.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title cannot be empty"})
			return
		}
		game.Title = *req.Title
	}
	if req.ReleaseYear != nil {
		game.ReleaseYear = *req.ReleaseYear
	}
	if req.PublisherID != nil {
		game.PublisherID = *req.PublisherID
	}
	// series_id の処理：キーが存在する場合は更新（null の場合は nil に設定）
	if _, exists := updateMap["series_id"]; exists {
		game.SeriesID = req.SeriesID
	}
	if req.Price != nil {
		game.Price = *req.Price
	}

	// PlatformIDsの更新
	if req.PlatformIDs != nil {
		var platforms []models.Platform
		if err := db.Where("id IN ?", req.PlatformIDs).Find(&platforms).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get platforms"})
			return
		}
		if len(platforms) != len(req.PlatformIDs) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "some platform_ids do not exist"})
			return
		}
		game.Platforms = platforms
	}

	// GenreIDsの更新
	if req.GenreIDs != nil {
		var genres []models.Genre
		if len(req.GenreIDs) > 0 {
			if err := db.Where("id IN ?", req.GenreIDs).Find(&genres).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get genres"})
				return
			}
			if len(genres) != len(req.GenreIDs) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "some genre_ids do not exist"})
				return
			}
		}
		game.Genres = genres
	}

	// バリデーション
	if err := utils.ValidateGame(game, db); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// データベースに保存（多対多リレーションも自動更新される）
	if err := db.Save(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update game"})
		return
	}

	// Preloadを使用して関連データを取得
	if err := db.Preload("Publisher").Preload("Platforms").Preload("Series").Preload("Genres").
		First(&game, game.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load game relations"})
		return
	}

	c.JSON(http.StatusOK, game)
}

// DeleteGame ゲーム削除ハンドラー（正規化後）
// GORMが自動的に中間テーブル（game_platforms, game_genres）のレコードも削除する
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

	// GORMが自動的に中間テーブルのレコードも削除する
	if err := db.Delete(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete game"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetStatistics 統計情報取得ハンドラー（正規化後）
func GetStatistics(c *gin.Context, db *gorm.DB) {
	var totalCount int64
	if err := db.Model(&models.Game{}).Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get total count"})
		return
	}

	// プラットフォーム別の集計（多対多リレーション、削除されたゲームを除外）
	type PlatformCount struct {
		Name  string
		Count int64
	}
	var platformCounts []PlatformCount
	if err := db.Model(&models.Platform{}).
		Select("platforms.name, COUNT(DISTINCT game_platforms.game_id) as count").
		Joins("LEFT JOIN game_platforms ON platforms.id = game_platforms.platform_id").
		Joins("INNER JOIN games ON game_platforms.game_id = games.id AND games.deleted_at IS NULL").
		Group("platforms.id, platforms.name").
		Scan(&platformCounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get platform counts"})
		return
	}
	platformCountsMap := make(map[string]int)
	for _, pc := range platformCounts {
		platformCountsMap[pc.Name] = int(pc.Count)
	}

	// 発売会社別の集計（1対多リレーション、削除されたゲームを除外）
	type PublisherCount struct {
		Name  string
		Count int64
	}
	var publisherCounts []PublisherCount
	if err := db.Model(&models.Publisher{}).
		Select("publishers.name, COUNT(games.id) as count").
		Joins("LEFT JOIN games ON publishers.id = games.publisher_id AND games.deleted_at IS NULL").
		Group("publishers.id, publishers.name").
		Scan(&publisherCounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get publisher counts"})
		return
	}
	publisherCountsMap := make(map[string]int)
	for _, pc := range publisherCounts {
		publisherCountsMap[pc.Name] = int(pc.Count)
	}

	// ジャンル別の集計（多対多リレーション、削除されたゲームを除外）
	type GenreCount struct {
		Name  string
		Count int64
	}
	var genreCounts []GenreCount
	if err := db.Model(&models.Genre{}).
		Select("genres.name, COUNT(DISTINCT game_genres.game_id) as count").
		Joins("LEFT JOIN game_genres ON genres.id = game_genres.genre_id").
		Joins("INNER JOIN games ON game_genres.game_id = games.id AND games.deleted_at IS NULL").
		Group("genres.id, genres.name").
		Scan(&genreCounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get genre counts"})
		return
	}
	genreCountsMap := make(map[string]int)
	for _, gc := range genreCounts {
		genreCountsMap[gc.Name] = int(gc.Count)
	}

	// シリーズ別の集計（1対多リレーション、削除されたゲームを除外）
	type SeriesCount struct {
		Name  string
		Count int64
	}
	var seriesCounts []SeriesCount
	if err := db.Model(&models.Series{}).
		Select("series.name, COUNT(games.id) as count").
		Joins("LEFT JOIN games ON series.id = games.series_id AND games.deleted_at IS NULL").
		Group("series.id, series.name").
		Scan(&seriesCounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get series counts"})
		return
	}
	seriesCountsMap := make(map[string]int)
	for _, sc := range seriesCounts {
		seriesCountsMap[sc.Name] = int(sc.Count)
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
