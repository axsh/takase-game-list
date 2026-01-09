package integration

import (
	"testing"

	"takase-game-list/models"
)

// TestNormalizedModelMigration_Integration 正規化モデルのマイグレーション確認
func TestNormalizedModelMigration_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	// Publisherテーブルの存在確認
	var publisher models.Publisher
	if !database.Migrator().HasTable(&publisher) {
		t.Error("publishers table does not exist")
	}

	// Platformテーブルの存在確認
	var platform models.Platform
	if !database.Migrator().HasTable(&platform) {
		t.Error("platforms table does not exist")
	}

	// Seriesテーブルの存在確認
	var series models.Series
	if !database.Migrator().HasTable(&series) {
		t.Error("series table does not exist")
	}

	// Genreテーブルの存在確認
	var genre models.Genre
	if !database.Migrator().HasTable(&genre) {
		t.Error("genres table does not exist")
	}

	// Gameテーブルの存在確認
	var game models.Game
	if !database.Migrator().HasTable(&game) {
		t.Error("games table does not exist")
	}

	// 中間テーブルの存在確認
	if !database.Migrator().HasTable("game_platforms") {
		t.Error("game_platforms table does not exist")
	}

	if !database.Migrator().HasTable("game_genres") {
		t.Error("game_genres table does not exist")
	}

	// Gameテーブルのカラム構造確認
	columnTypes, err := database.Migrator().ColumnTypes(&game)
	if err != nil {
		t.Fatalf("failed to get column types: %v", err)
	}

	// 必須カラムの確認（正規化後）
	// 注意: Platformsは多対多リレーションのため、game_platformsテーブルに保存される
	requiredColumns := map[string]bool{
		"id":           false,
		"title":        false,
		"release_year": false,
		"publisher_id": false,
		"series_id":    false,
		"price":        false,
		"created_at":   false,
		"updated_at":   false,
	}

	for _, ct := range columnTypes {
		columnName := ct.Name()
		if _, exists := requiredColumns[columnName]; exists {
			requiredColumns[columnName] = true
		}
	}

	for column, exists := range requiredColumns {
		if !exists {
			t.Errorf("required column '%s' does not exist", column)
		}
	}
}

// TestPublisherGameRelation_Integration PublisherとGameの1対多関係の確認
func TestPublisherGameRelation_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	// Publisherの作成
	publisher := models.Publisher{
		Name: "Test Publisher",
	}
	if err := database.Create(&publisher).Error; err != nil {
		t.Fatalf("failed to create publisher: %v", err)
	}

	// Gameの作成（PublisherIDを指定）
	game1 := models.Game{
		Title:       "Game 1",
		ReleaseYear: 2024,
		PublisherID: publisher.ID,
	}
	if err := database.Create(&game1).Error; err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	game2 := models.Game{
		Title:       "Game 2",
		ReleaseYear: 2023,
		PublisherID: publisher.ID,
	}
	if err := database.Create(&game2).Error; err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	// PublisherからGameを取得（Preload）
	var retrievedPublisher models.Publisher
	if err := database.Preload("Games").First(&retrievedPublisher, publisher.ID).Error; err != nil {
		t.Fatalf("failed to retrieve publisher: %v", err)
	}

	if len(retrievedPublisher.Games) != 2 {
		t.Errorf("expected 2 games, got %d", len(retrievedPublisher.Games))
	}

	// GameからPublisherを取得（Preload）
	var retrievedGame models.Game
	if err := database.Preload("Publisher").First(&retrievedGame, game1.ID).Error; err != nil {
		t.Fatalf("failed to retrieve game: %v", err)
	}

	if retrievedGame.Publisher.ID != publisher.ID {
		t.Errorf("expected publisher ID %d, got %d", publisher.ID, retrievedGame.Publisher.ID)
	}
	if retrievedGame.Publisher.Name != "Test Publisher" {
		t.Errorf("expected publisher name 'Test Publisher', got '%s'", retrievedGame.Publisher.Name)
	}
}

// TestPlatformGameRelation_Integration PlatformとGameの多対多関係の確認
func TestPlatformGameRelation_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	// Publisherの作成（必須）
	publisher := models.Publisher{
		Name: "Test Publisher",
	}
	if err := database.Create(&publisher).Error; err != nil {
		t.Fatalf("failed to create publisher: %v", err)
	}

	// Platformの作成
	platform1 := models.Platform{
		Name: "PC",
	}
	if err := database.Create(&platform1).Error; err != nil {
		t.Fatalf("failed to create platform: %v", err)
	}

	platform2 := models.Platform{
		Name: "Nintendo Switch",
	}
	if err := database.Create(&platform2).Error; err != nil {
		t.Fatalf("failed to create platform: %v", err)
	}

	// Gameの作成（複数のPlatformを関連付け）
	game := models.Game{
		Title:       "Multi-Platform Game",
		ReleaseYear: 2024,
		PublisherID: publisher.ID,
		Platforms:   []models.Platform{platform1, platform2},
	}
	if err := database.Create(&game).Error; err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	// GameからPlatformを取得（Preload）
	var retrievedGame models.Game
	if err := database.Preload("Platforms").First(&retrievedGame, game.ID).Error; err != nil {
		t.Fatalf("failed to retrieve game: %v", err)
	}

	if len(retrievedGame.Platforms) != 2 {
		t.Errorf("expected 2 platforms, got %d", len(retrievedGame.Platforms))
	}

	// PlatformからGameを取得（Preload）
	var retrievedPlatform models.Platform
	if err := database.Preload("Games").First(&retrievedPlatform, platform1.ID).Error; err != nil {
		t.Fatalf("failed to retrieve platform: %v", err)
	}

	if len(retrievedPlatform.Games) != 1 {
		t.Errorf("expected 1 game, got %d", len(retrievedPlatform.Games))
	}
}

// TestSeriesGameRelation_Integration SeriesとGameの1対多関係（NULL許容）の確認
func TestSeriesGameRelation_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	// Publisherの作成（必須）
	publisher := models.Publisher{
		Name: "Test Publisher",
	}
	if err := database.Create(&publisher).Error; err != nil {
		t.Fatalf("failed to create publisher: %v", err)
	}

	// Platformの作成（必須）
	platform := models.Platform{
		Name: "PC",
	}
	if err := database.Create(&platform).Error; err != nil {
		t.Fatalf("failed to create platform: %v", err)
	}

	// Seriesの作成
	series := models.Series{
		Name: "Test Series",
	}
	if err := database.Create(&series).Error; err != nil {
		t.Fatalf("failed to create series: %v", err)
	}

	// Gameの作成（SeriesIDを指定）
	game1 := models.Game{
		Title:       "Game 1",
		ReleaseYear: 2024,
		PublisherID: publisher.ID,
		SeriesID:    &series.ID,
		Platforms:   []models.Platform{platform},
	}
	if err := database.Create(&game1).Error; err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	// Gameの作成（SeriesIDをnil）
	game2 := models.Game{
		Title:       "Game 2",
		ReleaseYear: 2023,
		PublisherID: publisher.ID,
		SeriesID:    nil,
		Platforms:   []models.Platform{platform},
	}
	if err := database.Create(&game2).Error; err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	// SeriesからGameを取得（Preload）
	var retrievedSeries models.Series
	if err := database.Preload("Games").First(&retrievedSeries, series.ID).Error; err != nil {
		t.Fatalf("failed to retrieve series: %v", err)
	}

	if len(retrievedSeries.Games) != 1 {
		t.Errorf("expected 1 game, got %d", len(retrievedSeries.Games))
	}

	// GameからSeriesを取得（Preload、nil許容）
	var retrievedGame1 models.Game
	if err := database.Preload("Series").First(&retrievedGame1, game1.ID).Error; err != nil {
		t.Fatalf("failed to retrieve game: %v", err)
	}

	if retrievedGame1.Series == nil {
		t.Error("expected series to be set, got nil")
	} else if retrievedGame1.Series.ID != series.ID {
		t.Errorf("expected series ID %d, got %d", series.ID, retrievedGame1.Series.ID)
	}

	var retrievedGame2 models.Game
	if err := database.Preload("Series").First(&retrievedGame2, game2.ID).Error; err != nil {
		t.Fatalf("failed to retrieve game: %v", err)
	}

	if retrievedGame2.Series != nil {
		t.Error("expected series to be nil, got non-nil")
	}
}

// TestGenreGameRelation_Integration GenreとGameの多対多関係の確認
func TestGenreGameRelation_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	// Publisherの作成（必須）
	publisher := models.Publisher{
		Name: "Test Publisher",
	}
	if err := database.Create(&publisher).Error; err != nil {
		t.Fatalf("failed to create publisher: %v", err)
	}

	// Platformの作成（必須）
	platform := models.Platform{
		Name: "PC",
	}
	if err := database.Create(&platform).Error; err != nil {
		t.Fatalf("failed to create platform: %v", err)
	}

	// Genreの作成
	genre1 := models.Genre{
		Name: "RPG",
	}
	if err := database.Create(&genre1).Error; err != nil {
		t.Fatalf("failed to create genre: %v", err)
	}

	genre2 := models.Genre{
		Name: "Action",
	}
	if err := database.Create(&genre2).Error; err != nil {
		t.Fatalf("failed to create genre: %v", err)
	}

	// Gameの作成（複数のGenreを関連付け）
	game := models.Game{
		Title:       "Multi-Genre Game",
		ReleaseYear: 2024,
		PublisherID: publisher.ID,
		Platforms:   []models.Platform{platform},
		Genres:      []models.Genre{genre1, genre2},
	}
	if err := database.Create(&game).Error; err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	// GameからGenreを取得（Preload）
	var retrievedGame models.Game
	if err := database.Preload("Genres").First(&retrievedGame, game.ID).Error; err != nil {
		t.Fatalf("failed to retrieve game: %v", err)
	}

	if len(retrievedGame.Genres) != 2 {
		t.Errorf("expected 2 genres, got %d", len(retrievedGame.Genres))
	}

	// GenreからGameを取得（Preload）
	var retrievedGenre models.Genre
	if err := database.Preload("Games").First(&retrievedGenre, genre1.ID).Error; err != nil {
		t.Fatalf("failed to retrieve genre: %v", err)
	}

	if len(retrievedGenre.Games) != 1 {
		t.Errorf("expected 1 game, got %d", len(retrievedGenre.Games))
	}

	// GenreなしのGameも作成できること
	gameWithoutGenre := models.Game{
		Title:       "Game Without Genre",
		ReleaseYear: 2023,
		PublisherID: publisher.ID,
		Platforms:   []models.Platform{platform},
		Genres:      []models.Genre{},
	}
	if err := database.Create(&gameWithoutGenre).Error; err != nil {
		t.Fatalf("failed to create game without genre: %v", err)
	}
}

// TestCreateGameWithRelations_Integration リレーションを含むGame作成の確認
func TestCreateGameWithRelations_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	// マスタデータの作成
	publisher := models.Publisher{Name: "Test Publisher"}
	if err := database.Create(&publisher).Error; err != nil {
		t.Fatalf("failed to create publisher: %v", err)
	}

	platform1 := models.Platform{Name: "PC"}
	platform2 := models.Platform{Name: "Nintendo Switch"}
	if err := database.Create(&platform1).Error; err != nil {
		t.Fatalf("failed to create platform: %v", err)
	}
	if err := database.Create(&platform2).Error; err != nil {
		t.Fatalf("failed to create platform: %v", err)
	}

	series := models.Series{Name: "Test Series"}
	if err := database.Create(&series).Error; err != nil {
		t.Fatalf("failed to create series: %v", err)
	}

	genre1 := models.Genre{Name: "RPG"}
	genre2 := models.Genre{Name: "Action"}
	if err := database.Create(&genre1).Error; err != nil {
		t.Fatalf("failed to create genre: %v", err)
	}
	if err := database.Create(&genre2).Error; err != nil {
		t.Fatalf("failed to create genre: %v", err)
	}

	// Gameの作成（全リレーションを含む）
	game := models.Game{
		Title:       "Complete Game",
		ReleaseYear: 2024,
		PublisherID: publisher.ID,
		SeriesID:    &series.ID,
		Platforms:   []models.Platform{platform1, platform2},
		Genres:      []models.Genre{genre1, genre2},
		Price:       5000,
	}
	if err := database.Create(&game).Error; err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	// 全リレーションをPreloadして取得
	var retrievedGame models.Game
	if err := database.Preload("Publisher").Preload("Platforms").Preload("Series").Preload("Genres").
		First(&retrievedGame, game.ID).Error; err != nil {
		t.Fatalf("failed to retrieve game: %v", err)
	}

	// 各リレーションの確認
	if retrievedGame.Publisher.ID != publisher.ID {
		t.Errorf("expected publisher ID %d, got %d", publisher.ID, retrievedGame.Publisher.ID)
	}

	if len(retrievedGame.Platforms) != 2 {
		t.Errorf("expected 2 platforms, got %d", len(retrievedGame.Platforms))
	}

	if retrievedGame.Series == nil || retrievedGame.Series.ID != series.ID {
		t.Errorf("expected series ID %d, got %v", series.ID, retrievedGame.Series)
	}

	if len(retrievedGame.Genres) != 2 {
		t.Errorf("expected 2 genres, got %d", len(retrievedGame.Genres))
	}
}

// TestGetGameWithPreload_Integration Preloadによる関連データ取得の確認
func TestGetGameWithPreload_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	// マスタデータの作成
	publisher := models.Publisher{Name: "Test Publisher"}
	platform := models.Platform{Name: "PC"}
	series := models.Series{Name: "Test Series"}
	genre := models.Genre{Name: "RPG"}

	database.Create(&publisher)
	database.Create(&platform)
	database.Create(&series)
	database.Create(&genre)

	// Gameの作成
	game := models.Game{
		Title:       "Test Game",
		ReleaseYear: 2024,
		PublisherID: publisher.ID,
		SeriesID:    &series.ID,
		Platforms:   []models.Platform{platform},
		Genres:      []models.Genre{genre},
	}
	database.Create(&game)

	// Preloadなしで取得
	var gameWithoutPreload models.Game
	if err := database.First(&gameWithoutPreload, game.ID).Error; err != nil {
		t.Fatalf("failed to retrieve game: %v", err)
	}

	// リレーションデータが空であることを確認
	if gameWithoutPreload.Publisher.ID != 0 {
		t.Error("expected publisher to be empty without preload")
	}
	if len(gameWithoutPreload.Platforms) != 0 {
		t.Error("expected platforms to be empty without preload")
	}

	// Preloadありで取得
	var gameWithPreload models.Game
	if err := database.Preload("Publisher").Preload("Platforms").Preload("Series").Preload("Genres").
		First(&gameWithPreload, game.ID).Error; err != nil {
		t.Fatalf("failed to retrieve game: %v", err)
	}

	// リレーションデータが設定されていることを確認
	if gameWithPreload.Publisher.ID == 0 {
		t.Error("expected publisher to be loaded with preload")
	}
	if len(gameWithPreload.Platforms) == 0 {
		t.Error("expected platforms to be loaded with preload")
	}
	if gameWithPreload.Series == nil {
		t.Error("expected series to be loaded with preload")
	}
	if len(gameWithPreload.Genres) == 0 {
		t.Error("expected genres to be loaded with preload")
	}
}

// TestUpdateGameRelations_Integration 多対多リレーションの更新確認
func TestUpdateGameRelations_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	// マスタデータの作成
	publisher := models.Publisher{Name: "Test Publisher"}
	platform1 := models.Platform{Name: "PC"}
	platform2 := models.Platform{Name: "Nintendo Switch"}
	platform3 := models.Platform{Name: "PlayStation 5"}

	database.Create(&publisher)
	database.Create(&platform1)
	database.Create(&platform2)
	database.Create(&platform3)

	// Gameの作成（platform1とplatform2を関連付け）
	game := models.Game{
		Title:       "Test Game",
		ReleaseYear: 2024,
		PublisherID: publisher.ID,
		Platforms:   []models.Platform{platform1, platform2},
	}
	database.Create(&game)

	// Platformを更新（platform2とplatform3に変更）
	var retrievedGame models.Game
	database.Preload("Platforms").First(&retrievedGame, game.ID)
	// GORMのAssociationを使用して多対多リレーションを更新
	if err := database.Model(&retrievedGame).Association("Platforms").Replace([]models.Platform{platform2, platform3}); err != nil {
		t.Fatalf("failed to update game platforms: %v", err)
	}

	// 更新後の確認
	var updatedGame models.Game
	database.Preload("Platforms").First(&updatedGame, game.ID)

	if len(updatedGame.Platforms) != 2 {
		t.Errorf("expected 2 platforms, got %d", len(updatedGame.Platforms))
	}

	platformIDs := make(map[uint]bool)
	for _, p := range updatedGame.Platforms {
		platformIDs[p.ID] = true
	}

	if !platformIDs[platform2.ID] || !platformIDs[platform3.ID] {
		t.Error("expected platforms to be updated to platform2 and platform3")
	}
}

// TestDeleteGameWithRelations_Integration 中間テーブルのレコード削除確認
func TestDeleteGameWithRelations_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	// マスタデータの作成
	publisher := models.Publisher{Name: "Test Publisher"}
	platform := models.Platform{Name: "PC"}
	genre := models.Genre{Name: "RPG"}

	database.Create(&publisher)
	database.Create(&platform)
	database.Create(&genre)

	// Gameの作成
	game := models.Game{
		Title:       "Test Game",
		ReleaseYear: 2024,
		PublisherID: publisher.ID,
		Platforms:   []models.Platform{platform},
		Genres:      []models.Genre{genre},
	}
	database.Create(&game)

	// 中間テーブルのレコード数を確認
	var count int64
	database.Table("game_platforms").Where("game_id = ?", game.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 record in game_platforms, got %d", count)
	}

	database.Table("game_genres").Where("game_id = ?", game.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 record in game_genres, got %d", count)
	}

	// Gameの削除前に中間テーブルのレコードを手動で削除
	// 注意: GORMのソフトデリートでは、中間テーブルのレコードは自動的に削除されません
	// そのため、Associationをクリアするか、手動で削除する必要があります
	database.Model(&game).Association("Platforms").Clear()
	database.Model(&game).Association("Genres").Clear()

	// Gameの削除（物理削除）
	if err := database.Unscoped().Delete(&game).Error; err != nil {
		t.Fatalf("failed to delete game: %v", err)
	}

	// 中間テーブルのレコードが削除されていることを確認
	database.Table("game_platforms").Where("game_id = ?", game.ID).Count(&count)
	if count != 0 {
		t.Errorf("expected 0 records in game_platforms after deletion, got %d", count)
	}

	database.Table("game_genres").Where("game_id = ?", game.ID).Count(&count)
	if count != 0 {
		t.Errorf("expected 0 records in game_genres after deletion, got %d", count)
	}
}
