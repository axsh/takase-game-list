package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"takase-game-list/handlers"
	"takase-game-list/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// setupTestServer テスト用のHTTPサーバーをセットアップする
func setupTestServer(t *testing.T, database *gorm.DB) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// ゲームAPIのルーティング
	router.POST("/games", func(c *gin.Context) {
		handlers.CreateGame(c, database)
	})
	router.GET("/games", func(c *gin.Context) {
		handlers.GetGames(c, database)
	})
	// GET /games/:id はフェーズ2では実装しない（フェーズ3以降で実装予定）
	router.GET("/games/search", func(c *gin.Context) {
		handlers.SearchGames(c, database)
	})
	router.GET("/games/statistics", func(c *gin.Context) {
		handlers.GetStatistics(c, database)
	})
	router.PUT("/games/:id", func(c *gin.Context) {
		handlers.UpdateGame(c, database)
	})
	router.DELETE("/games/:id", func(c *gin.Context) {
		handlers.DeleteGame(c, database)
	})

	return router
}

// TestCreateGame_Integration ゲーム登録の結合テスト
func TestCreateGame_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServer(t, database)

	t.Run("必須項目のみで登録できること", func(t *testing.T) {
		// マスタデータの作成
		publisher := models.Publisher{Name: "Test Publisher"}
		database.Create(&publisher)

		platform := models.Platform{Name: "PC"}
		database.Create(&platform)

		reqBody := map[string]interface{}{
			"title":        "Test Game",
			"release_year": 2024,
			"publisher_id": publisher.ID,
			"platform_ids": []uint{platform.ID},
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var game models.Game
		if err := json.Unmarshal(w.Body.Bytes(), &game); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if game.ID == 0 {
			t.Error("game ID was not auto-generated")
		}
		if game.Title != "Test Game" {
			t.Errorf("expected title 'Test Game', got '%s'", game.Title)
		}
		if game.ReleaseYear != 2024 {
			t.Errorf("expected release year 2024, got %d", game.ReleaseYear)
		}
	})

	t.Run("全項目を指定して登録できること", func(t *testing.T) {
		// マスタデータの作成
		publisher := models.Publisher{Name: "Test Publisher 2"}
		database.Create(&publisher)

		platform := models.Platform{Name: "Nintendo Switch"}
		database.Create(&platform)

		series := models.Series{Name: "Test Series"}
		database.Create(&series)

		genre := models.Genre{Name: "RPG"}
		database.Create(&genre)

		reqBody := map[string]interface{}{
			"title":        "Test Game 2",
			"release_year": 2023,
			"publisher_id": publisher.ID,
			"platform_ids": []uint{platform.ID},
			"series_id":    series.ID,
			"genre_ids":    []uint{genre.ID},
			"price":        5000,
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var game models.Game
		if err := json.Unmarshal(w.Body.Bytes(), &game); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if game.Series == nil || game.Series.Name != "Test Series" {
			t.Errorf("expected series 'Test Series', got %v", game.Series)
		}
		if len(game.Genres) != 1 || game.Genres[0].Name != "RPG" {
			t.Errorf("expected genre 'RPG', got %v", game.Genres)
		}
		if game.Price != 5000 {
			t.Errorf("expected price 5000, got %d", game.Price)
		}
	})

	t.Run("必須項目不足時の400エラー", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title": "Test Game",
			// release_year, publisher, platform が不足
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

// TestGetGames_Integration ゲーム一覧取得の結合テスト
func TestGetGames_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServer(t, database)

	// マスタデータの作成
	publisher1 := models.Publisher{Name: "Publisher 1"}
	publisher2 := models.Publisher{Name: "Publisher 2"}
	database.Create(&publisher1)
	database.Create(&publisher2)

	platform1 := models.Platform{Name: "PC"}
	platform2 := models.Platform{Name: "Nintendo Switch"}
	database.Create(&platform1)
	database.Create(&platform2)

	// テストデータの準備
	games := []models.Game{
		{Title: "Game 1", ReleaseYear: 2024, PublisherID: publisher1.ID, Platforms: []models.Platform{platform1}},
		{Title: "Game 2", ReleaseYear: 2023, PublisherID: publisher2.ID, Platforms: []models.Platform{platform2}},
	}
	for i := range games {
		database.Create(&games[i])
	}

	t.Run("全件取得できること", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/games", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var result []models.Game
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("expected 2 games, got %d", len(result))
		}
	})

	t.Run("空の場合は空配列が返ること", func(t *testing.T) {
		emptyDB, emptyPath := setupTestDB(t)
		defer cleanupTestDB(t, emptyDB, emptyPath)

		emptyRouter := setupTestServer(t, emptyDB)

		req := httptest.NewRequest(http.MethodGet, "/games", nil)
		w := httptest.NewRecorder()

		emptyRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var result []models.Game
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(result) != 0 {
			t.Errorf("expected 0 games, got %d", len(result))
		}
	})

	t.Run("ソートが正しく動作すること", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/games?sort=title&order=asc", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var result []models.Game
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(result) < 2 {
			t.Fatal("expected at least 2 games")
		}
		if result[0].Title > result[1].Title {
			t.Error("games are not sorted by title in ascending order")
		}
	})
}

// TestUpdateGame_Integration ゲーム更新の結合テスト
func TestUpdateGame_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServer(t, database)

	// マスタデータの作成
	publisher := models.Publisher{Name: "Original Publisher"}
	database.Create(&publisher)

	platform := models.Platform{Name: "PC"}
	database.Create(&platform)

	// テストデータの準備
	game := models.Game{
		Title:       "Original Game",
		ReleaseYear: 2024,
		PublisherID: publisher.ID,
		Platforms:   []models.Platform{platform},
	}
	database.Create(&game)

	t.Run("部分更新ができること", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title": "Updated Game",
		}
		jsonBody, _ := json.Marshal(reqBody)

		url := fmt.Sprintf("/games/%d", game.ID)
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var updatedGame models.Game
		if err := json.Unmarshal(w.Body.Bytes(), &updatedGame); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if updatedGame.Title != "Updated Game" {
			t.Errorf("expected title 'Updated Game', got '%s'", updatedGame.Title)
		}
		if updatedGame.Publisher.Name != "Original Publisher" {
			t.Errorf("expected publisher to remain 'Original Publisher', got '%s'", updatedGame.Publisher.Name)
		}
	})

	t.Run("存在しないIDの404エラー", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title": "Updated Game",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/games/99999", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("必須項目を空に更新しようとした場合の400エラー", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title": "",
		}
		jsonBody, _ := json.Marshal(reqBody)

		url := fmt.Sprintf("/games/%d", game.ID)
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

// TestDeleteGame_Integration ゲーム削除の結合テスト
func TestDeleteGame_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServer(t, database)

	// マスタデータの作成
	publisher := models.Publisher{Name: "Publisher"}
	database.Create(&publisher)

	platform := models.Platform{Name: "PC"}
	database.Create(&platform)

	// テストデータの準備
	game := models.Game{
		Title:       "Game to Delete",
		ReleaseYear: 2024,
		PublisherID: publisher.ID,
		Platforms:   []models.Platform{platform},
	}
	database.Create(&game)

	t.Run("存在するIDで削除できること", func(t *testing.T) {
		url := fmt.Sprintf("/games/%d", game.ID)
		req := httptest.NewRequest(http.MethodDelete, url, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
		}
	})

	t.Run("削除後に取得できないこと", func(t *testing.T) {
		// 削除されたゲームが一覧に含まれないことを確認
		req := httptest.NewRequest(http.MethodGet, "/games", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var games []models.Game
		if err := json.Unmarshal(w.Body.Bytes(), &games); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		// 削除されたゲームが一覧に含まれていないことを確認
		for _, g := range games {
			if g.ID == game.ID {
				t.Errorf("deleted game with ID %d should not be in the list", game.ID)
			}
		}
	})

	t.Run("存在しないIDの404エラー", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/games/99999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

// TestSearchGames_Integration 検索機能の結合テスト
func TestSearchGames_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServer(t, database)

	// マスタデータの作成
	publisher1 := models.Publisher{Name: "Square Enix"}
	publisher2 := models.Publisher{Name: "Nintendo"}
	database.Create(&publisher1)
	database.Create(&publisher2)

	platform1 := models.Platform{Name: "PC"}
	platform2 := models.Platform{Name: "Nintendo Switch"}
	database.Create(&platform1)
	database.Create(&platform2)

	// テストデータの準備
	games := []models.Game{
		{Title: "Final Fantasy VII", ReleaseYear: 2020, PublisherID: publisher1.ID, Platforms: []models.Platform{platform1}},
		{Title: "Final Fantasy XV", ReleaseYear: 2016, PublisherID: publisher1.ID, Platforms: []models.Platform{platform1}},
		{Title: "The Legend of Zelda", ReleaseYear: 2017, PublisherID: publisher2.ID, Platforms: []models.Platform{platform2}},
		{Title: "Super Mario Odyssey", ReleaseYear: 2017, PublisherID: publisher2.ID, Platforms: []models.Platform{platform2}},
	}
	for i := range games {
		database.Create(&games[i])
	}

	t.Run("タイトル部分一致検索が動作すること", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/games/search?q=Final", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var result []models.Game
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("expected 2 games, got %d", len(result))
		}
		for _, game := range result {
			if !contains(game.Title, "Final") {
				t.Errorf("game title '%s' does not contain 'Final'", game.Title)
			}
		}
	})

	t.Run("大文字・小文字を区別しないこと", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/games/search?q=final", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var result []models.Game
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("expected 2 games, got %d", len(result))
		}
	})

	t.Run("検索結果が0件の場合、空配列が返ること", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/games/search?q=Nonexistent", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var result []models.Game
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(result) != 0 {
			t.Errorf("expected 0 games, got %d", len(result))
		}
	})
}

// TestFilterGames_Integration フィルタリング機能の結合テスト
func TestFilterGames_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServer(t, database)

	// マスタデータの作成
	publisher1 := models.Publisher{Name: "Square Enix"}
	publisher2 := models.Publisher{Name: "Nintendo"}
	database.Create(&publisher1)
	database.Create(&publisher2)

	platform1 := models.Platform{Name: "PC"}
	platform2 := models.Platform{Name: "Nintendo Switch"}
	database.Create(&platform1)
	database.Create(&platform2)

	series1 := models.Series{Name: "Final Fantasy"}
	series2 := models.Series{Name: "The Legend of Zelda"}
	database.Create(&series1)
	database.Create(&series2)

	genre1 := models.Genre{Name: "RPG"}
	genre2 := models.Genre{Name: "Action"}
	database.Create(&genre1)
	database.Create(&genre2)

	// テストデータの準備
	games := []models.Game{
		{Title: "Final Fantasy VII", ReleaseYear: 2020, PublisherID: publisher1.ID, SeriesID: &series1.ID, Platforms: []models.Platform{platform1}, Genres: []models.Genre{genre1}},
		{Title: "Final Fantasy XV", ReleaseYear: 2016, PublisherID: publisher1.ID, SeriesID: &series1.ID, Platforms: []models.Platform{platform1}, Genres: []models.Genre{genre1}},
		{Title: "The Legend of Zelda", ReleaseYear: 2017, PublisherID: publisher2.ID, SeriesID: &series2.ID, Platforms: []models.Platform{platform2}, Genres: []models.Genre{genre2}},
		{Title: "Super Mario Odyssey", ReleaseYear: 2017, PublisherID: publisher2.ID, Platforms: []models.Platform{platform2}, Genres: []models.Genre{genre2}},
	}
	for i := range games {
		database.Create(&games[i])
	}

	t.Run("プラットフォームフィルタが動作すること", func(t *testing.T) {
		// platform_idsでフィルタリング（正規化後）
		var platform models.Platform
		database.Where("name = ?", "PC").First(&platform)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/games?platform_ids=%d", platform.ID), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var result []models.Game
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("expected 2 games, got %d", len(result))
		}
		for _, game := range result {
			if len(game.Platforms) == 0 || game.Platforms[0].Name != "PC" {
				t.Errorf("expected platform 'PC', got %v", game.Platforms)
			}
		}
	})

	t.Run("発売会社フィルタが動作すること", func(t *testing.T) {
		// publisher_idでフィルタリング（正規化後）
		var publisher models.Publisher
		database.Where("name = ?", "Square Enix").First(&publisher)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/games?publisher_id=%d", publisher.ID), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var result []models.Game
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("expected 2 games, got %d", len(result))
		}
		for _, game := range result {
			if game.Publisher.Name != "Square Enix" {
				t.Errorf("expected publisher 'Square Enix', got '%s'", game.Publisher.Name)
			}
		}
	})

	t.Run("ジャンルフィルタが動作すること", func(t *testing.T) {
		// genre_idsでフィルタリング（正規化後）
		var genre models.Genre
		database.Where("name = ?", "RPG").First(&genre)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/games?genre_ids=%d", genre.ID), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var result []models.Game
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("expected 2 games, got %d", len(result))
		}
		for _, game := range result {
			if len(game.Genres) == 0 || game.Genres[0].Name != "RPG" {
				t.Errorf("expected genre 'RPG', got %v", game.Genres)
			}
		}
	})

	t.Run("シリーズフィルタが動作すること", func(t *testing.T) {
		// series_idでフィルタリング（正規化後）
		var series models.Series
		database.Where("name = ?", "Final Fantasy").First(&series)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/games?series_id=%d", series.ID), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var result []models.Game
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("expected 2 games, got %d", len(result))
		}
		for _, game := range result {
			if game.Series == nil || game.Series.Name != "Final Fantasy" {
				t.Errorf("expected series 'Final Fantasy', got %v", game.Series)
			}
		}
	})

	t.Run("発売年範囲フィルタが動作すること", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/games?min_year=2017&max_year=2020", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var result []models.Game
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(result) != 3 {
			t.Errorf("expected 3 games, got %d", len(result))
		}
		for _, game := range result {
			if game.ReleaseYear < 2017 || game.ReleaseYear > 2020 {
				t.Errorf("expected release year between 2017 and 2020, got %d", game.ReleaseYear)
			}
		}
	})

	t.Run("複数条件の組み合わせが動作すること", func(t *testing.T) {
		// platform_idsとpublisher_idでフィルタリング（正規化後）
		var platform models.Platform
		var publisher models.Publisher
		database.Where("name = ?", "Nintendo Switch").First(&platform)
		database.Where("name = ?", "Nintendo").First(&publisher)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/games?platform_ids=%d&publisher_id=%d", platform.ID, publisher.ID), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var result []models.Game
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("expected 2 games, got %d", len(result))
		}
		for _, game := range result {
			if len(game.Platforms) == 0 || game.Platforms[0].Name != "Nintendo Switch" {
				t.Errorf("expected platform 'Nintendo Switch', got %v", game.Platforms)
			}
			if game.Publisher.Name != "Nintendo" {
				t.Errorf("expected publisher 'Nintendo', got '%s'", game.Publisher.Name)
			}
		}
	})
}

// contains 文字列が部分文字列を含むかどうかをチェック（大文字・小文字を区別しない）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		containsIgnoreCase(s, substr))
}

func containsIgnoreCase(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

// TestGetStatistics_Integration 統計情報取得の結合テスト
func TestGetStatistics_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServer(t, database)

	// マスタデータの作成
	publisher1 := models.Publisher{Name: "Square Enix"}
	publisher2 := models.Publisher{Name: "Nintendo"}
	publisher3 := models.Publisher{Name: "Publisher"}
	database.Create(&publisher1)
	database.Create(&publisher2)
	database.Create(&publisher3)

	platform1 := models.Platform{Name: "PC"}
	platform2 := models.Platform{Name: "Nintendo Switch"}
	database.Create(&platform1)
	database.Create(&platform2)

	series1 := models.Series{Name: "Final Fantasy"}
	series2 := models.Series{Name: "The Legend of Zelda"}
	database.Create(&series1)
	database.Create(&series2)

	genre1 := models.Genre{Name: "RPG"}
	genre2 := models.Genre{Name: "Action"}
	database.Create(&genre1)
	database.Create(&genre2)

	// テストデータの準備
	games := []models.Game{
		{Title: "Final Fantasy VII", ReleaseYear: 2020, PublisherID: publisher1.ID, SeriesID: &series1.ID, Platforms: []models.Platform{platform1}, Genres: []models.Genre{genre1}, Price: 5000},
		{Title: "Final Fantasy XV", ReleaseYear: 2016, PublisherID: publisher1.ID, SeriesID: &series1.ID, Platforms: []models.Platform{platform1}, Genres: []models.Genre{genre1}, Price: 6000},
		{Title: "The Legend of Zelda", ReleaseYear: 2017, PublisherID: publisher2.ID, SeriesID: &series2.ID, Platforms: []models.Platform{platform2}, Genres: []models.Genre{genre2}, Price: 7000},
		{Title: "Super Mario Odyssey", ReleaseYear: 2017, PublisherID: publisher2.ID, Platforms: []models.Platform{platform2}, Genres: []models.Genre{genre2}, Price: 0}, // 価格0は除外
		{Title: "Game Without Genre", ReleaseYear: 2021, PublisherID: publisher3.ID, Platforms: []models.Platform{platform1}, Genres: []models.Genre{}, Price: 3000},   // ジャンルなし
	}
	for i := range games {
		database.Create(&games[i])
	}

	req := httptest.NewRequest(http.MethodGet, "/games/statistics", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	t.Run("総数が正しく計算されること", func(t *testing.T) {
		totalCount, ok := result["total_count"].(float64)
		if !ok {
			t.Fatal("total_count is not a number")
		}
		if int(totalCount) != 5 {
			t.Errorf("expected total_count 5, got %d", int(totalCount))
		}
	})

	t.Run("プラットフォーム別の集計が正しいこと", func(t *testing.T) {
		platformCounts, ok := result["platform_counts"].(map[string]interface{})
		if !ok {
			t.Fatal("platform_counts is not a map")
		}

		pcCount, ok := platformCounts["PC"].(float64)
		if !ok || int(pcCount) != 3 {
			t.Errorf("expected PC count 3, got %v", platformCounts["PC"])
		}

		switchCount, ok := platformCounts["Nintendo Switch"].(float64)
		if !ok || int(switchCount) != 2 {
			t.Errorf("expected Nintendo Switch count 2, got %v", platformCounts["Nintendo Switch"])
		}
	})

	t.Run("発売会社別の集計が正しいこと", func(t *testing.T) {
		publisherCounts, ok := result["publisher_counts"].(map[string]interface{})
		if !ok {
			t.Fatal("publisher_counts is not a map")
		}

		squareEnixCount, ok := publisherCounts["Square Enix"].(float64)
		if !ok || int(squareEnixCount) != 2 {
			t.Errorf("expected Square Enix count 2, got %v", publisherCounts["Square Enix"])
		}

		nintendoCount, ok := publisherCounts["Nintendo"].(float64)
		if !ok || int(nintendoCount) != 2 {
			t.Errorf("expected Nintendo count 2, got %v", publisherCounts["Nintendo"])
		}
	})

	t.Run("ジャンル別の集計が正しいこと（null除外）", func(t *testing.T) {
		genreCounts, ok := result["genre_counts"].(map[string]interface{})
		if !ok {
			t.Fatal("genre_counts is not a map")
		}

		rpgCount, ok := genreCounts["RPG"].(float64)
		if !ok || int(rpgCount) != 2 {
			t.Errorf("expected RPG count 2, got %v", genreCounts["RPG"])
		}

		actionCount, ok := genreCounts["Action"].(float64)
		if !ok || int(actionCount) != 2 {
			t.Errorf("expected Action count 2, got %v", genreCounts["Action"])
		}

		// ジャンルなしのゲームは集計に含まれない
		if _, exists := genreCounts[""]; exists {
			t.Error("genre_counts should not include empty genre")
		}
	})

	t.Run("シリーズ別の集計が正しいこと（null除外）", func(t *testing.T) {
		seriesCounts, ok := result["series_counts"].(map[string]interface{})
		if !ok {
			t.Fatal("series_counts is not a map")
		}

		ffCount, ok := seriesCounts["Final Fantasy"].(float64)
		if !ok || int(ffCount) != 2 {
			t.Errorf("expected Final Fantasy count 2, got %v", seriesCounts["Final Fantasy"])
		}

		zeldaCount, ok := seriesCounts["The Legend of Zelda"].(float64)
		if !ok || int(zeldaCount) != 1 {
			t.Errorf("expected The Legend of Zelda count 1, got %v", seriesCounts["The Legend of Zelda"])
		}

		// シリーズなしのゲームは集計に含まれない
		if _, exists := seriesCounts[""]; exists {
			t.Error("series_counts should not include empty series")
		}
	})

	t.Run("総購入金額が正しく計算されること（価格0の除外）", func(t *testing.T) {
		totalPrice, ok := result["total_price"].(float64)
		if !ok {
			t.Fatal("total_price is not a number")
		}
		// 5000 + 6000 + 7000 + 3000 = 21000（価格0のゲームは除外）
		expected := 21000
		if int(totalPrice) != expected {
			t.Errorf("expected total_price %d, got %d", expected, int(totalPrice))
		}
	})

	t.Run("平均価格が正しく計算されること（価格0の除外）", func(t *testing.T) {
		averagePrice, ok := result["average_price"].(float64)
		if !ok {
			t.Fatal("average_price is not a number")
		}
		// (5000 + 6000 + 7000 + 3000) / 4 = 5250
		expected := 5250.0
		if averagePrice != expected {
			t.Errorf("expected average_price %.2f, got %.2f", expected, averagePrice)
		}
	})
}

// TestGetStatistics_ExcludeDeletedGames_Integration 削除されたゲームが統計情報に含まれないことを確認するテスト
func TestGetStatistics_ExcludeDeletedGames_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServer(t, database)

	// マスタデータの作成
	publisher := models.Publisher{Name: "Test Publisher"}
	platform := models.Platform{Name: "PC"}
	series := models.Series{Name: "Test Series"}
	genre := models.Genre{Name: "RPG"}

	database.Create(&publisher)
	database.Create(&platform)
	database.Create(&series)
	database.Create(&genre)

	// テストデータの準備（削除前の統計）
	game1 := models.Game{
		Title:       "Game 1",
		ReleaseYear: 2020,
		PublisherID: publisher.ID,
		SeriesID:    &series.ID,
		Platforms:   []models.Platform{platform},
		Genres:      []models.Genre{genre},
		Price:       5000,
	}
	game2 := models.Game{
		Title:       "Game 2",
		ReleaseYear: 2021,
		PublisherID: publisher.ID,
		Platforms:   []models.Platform{platform},
		Genres:      []models.Genre{genre},
		Price:       6000,
	}
	database.Create(&game1)
	database.Create(&game2)

	// 削除前の統計情報を取得
	req1 := httptest.NewRequest(http.MethodGet, "/games/statistics", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w1.Code)
	}

	var result1 map[string]interface{}
	if err := json.Unmarshal(w1.Body.Bytes(), &result1); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// 削除前の確認
	totalCount1, _ := result1["total_count"].(float64)
	if int(totalCount1) != 2 {
		t.Errorf("expected total_count 2 before deletion, got %d", int(totalCount1))
	}

	// game1を削除（ソフトデリート）
	database.Delete(&game1)

	// 削除後の統計情報を取得
	req2 := httptest.NewRequest(http.MethodGet, "/games/statistics", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w2.Code)
	}

	var result2 map[string]interface{}
	if err := json.Unmarshal(w2.Body.Bytes(), &result2); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// 削除後の確認
	t.Run("総数が削除されたゲームを除外すること", func(t *testing.T) {
		totalCount2, ok := result2["total_count"].(float64)
		if !ok {
			t.Fatal("total_count is not a number")
		}
		if int(totalCount2) != 1 {
			t.Errorf("expected total_count 1 after deletion, got %d", int(totalCount2))
		}
	})

	t.Run("プラットフォーム別の集計が削除されたゲームを除外すること", func(t *testing.T) {
		platformCounts, ok := result2["platform_counts"].(map[string]interface{})
		if !ok {
			t.Fatal("platform_counts is not a map")
		}
		pcCount, ok := platformCounts["PC"].(float64)
		if !ok || int(pcCount) != 1 {
			t.Errorf("expected PC count 1 after deletion, got %v", platformCounts["PC"])
		}
	})

	t.Run("発売会社別の集計が削除されたゲームを除外すること", func(t *testing.T) {
		publisherCounts, ok := result2["publisher_counts"].(map[string]interface{})
		if !ok {
			t.Fatal("publisher_counts is not a map")
		}
		pubCount, ok := publisherCounts["Test Publisher"].(float64)
		if !ok || int(pubCount) != 1 {
			t.Errorf("expected Test Publisher count 1 after deletion, got %v", publisherCounts["Test Publisher"])
		}
	})

	t.Run("ジャンル別の集計が削除されたゲームを除外すること", func(t *testing.T) {
		genreCounts, ok := result2["genre_counts"].(map[string]interface{})
		if !ok {
			t.Fatal("genre_counts is not a map")
		}
		rpgCount, ok := genreCounts["RPG"].(float64)
		if !ok || int(rpgCount) != 1 {
			t.Errorf("expected RPG count 1 after deletion, got %v", genreCounts["RPG"])
		}
	})

	t.Run("シリーズ別の集計が削除されたゲームを除外すること", func(t *testing.T) {
		seriesCounts, ok := result2["series_counts"].(map[string]interface{})
		if !ok {
			t.Fatal("series_counts is not a map")
		}
		// シリーズが設定されていたgame1が削除されたので、シリーズ別の集計には含まれない
		testSeriesCount, exists := seriesCounts["Test Series"]
		if exists {
			count, ok := testSeriesCount.(float64)
			if ok && int(count) != 0 {
				t.Errorf("expected Test Series count 0 after deletion, got %v", testSeriesCount)
			}
		}
	})

	t.Run("発売年別の集計が削除されたゲームを除外すること", func(t *testing.T) {
		yearCounts, ok := result2["year_counts"].(map[string]interface{})
		if !ok {
			t.Fatal("year_counts is not a map")
		}
		year2020Count, exists := yearCounts["2020"]
		if exists {
			count, ok := year2020Count.(float64)
			if ok && int(count) != 0 {
				t.Errorf("expected year 2020 count 0 after deletion, got %v", year2020Count)
			}
		}
		year2021Count, ok := yearCounts["2021"].(float64)
		if !ok || int(year2021Count) != 1 {
			t.Errorf("expected year 2021 count 1 after deletion, got %v", yearCounts["2021"])
		}
	})

	t.Run("価格統計が削除されたゲームを除外すること", func(t *testing.T) {
		totalPrice, ok := result2["total_price"].(float64)
		if !ok {
			t.Fatal("total_price is not a number")
		}
		// game1(5000)が削除されたので、game2(6000)のみ
		expected := 6000
		if int(totalPrice) != expected {
			t.Errorf("expected total_price %d after deletion, got %d", expected, int(totalPrice))
		}

		averagePrice, ok := result2["average_price"].(float64)
		if !ok {
			t.Fatal("average_price is not a number")
		}
		expectedAvg := 6000.0
		if averagePrice != expectedAvg {
			t.Errorf("expected average_price %.2f after deletion, got %.2f", expectedAvg, averagePrice)
		}
	})
}
