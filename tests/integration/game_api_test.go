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
		reqBody := map[string]interface{}{
			"title":        "Test Game",
			"release_year": 2024,
			"publisher":    "Test Publisher",
			"platform":     "PC",
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
		series := "Test Series"
		genre := "RPG"
		reqBody := map[string]interface{}{
			"title":        "Test Game 2",
			"release_year": 2023,
			"publisher":    "Test Publisher 2",
			"platform":     "Nintendo Switch",
			"series":       series,
			"genre":        genre,
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

		if game.Series == nil || *game.Series != series {
			t.Errorf("expected series '%s', got %v", series, game.Series)
		}
		if game.Genre == nil || *game.Genre != genre {
			t.Errorf("expected genre '%s', got %v", genre, game.Genre)
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

	// テストデータの準備
	games := []models.Game{
		{Title: "Game 1", ReleaseYear: 2024, Publisher: "Publisher 1", Platform: "PC"},
		{Title: "Game 2", ReleaseYear: 2023, Publisher: "Publisher 2", Platform: "Nintendo Switch"},
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

	// テストデータの準備
	game := models.Game{
		Title:       "Original Game",
		ReleaseYear: 2024,
		Publisher:   "Original Publisher",
		Platform:    "PC",
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
		if updatedGame.Publisher != "Original Publisher" {
			t.Errorf("expected publisher to remain 'Original Publisher', got '%s'", updatedGame.Publisher)
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

	// テストデータの準備
	game := models.Game{
		Title:       "Game to Delete",
		ReleaseYear: 2024,
		Publisher:   "Publisher",
		Platform:    "PC",
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

	// テストデータの準備
	games := []models.Game{
		{Title: "Final Fantasy VII", ReleaseYear: 2020, Publisher: "Square Enix", Platform: "PC"},
		{Title: "Final Fantasy XV", ReleaseYear: 2016, Publisher: "Square Enix", Platform: "PC"},
		{Title: "The Legend of Zelda", ReleaseYear: 2017, Publisher: "Nintendo", Platform: "Nintendo Switch"},
		{Title: "Super Mario Odyssey", ReleaseYear: 2017, Publisher: "Nintendo", Platform: "Nintendo Switch"},
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

	series1 := "Final Fantasy"
	series2 := "The Legend of Zelda"
	genre1 := "RPG"
	genre2 := "Action"

	// テストデータの準備
	games := []models.Game{
		{Title: "Final Fantasy VII", ReleaseYear: 2020, Publisher: "Square Enix", Platform: "PC", Series: &series1, Genre: &genre1},
		{Title: "Final Fantasy XV", ReleaseYear: 2016, Publisher: "Square Enix", Platform: "PC", Series: &series1, Genre: &genre1},
		{Title: "The Legend of Zelda", ReleaseYear: 2017, Publisher: "Nintendo", Platform: "Nintendo Switch", Series: &series2, Genre: &genre2},
		{Title: "Super Mario Odyssey", ReleaseYear: 2017, Publisher: "Nintendo", Platform: "Nintendo Switch", Genre: &genre2},
	}
	for i := range games {
		database.Create(&games[i])
	}

	t.Run("プラットフォームフィルタが動作すること", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/games?platform=PC", nil)
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
			if game.Platform != "PC" {
				t.Errorf("expected platform 'PC', got '%s'", game.Platform)
			}
		}
	})

	t.Run("発売会社フィルタが動作すること", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/games?publisher=Square+Enix", nil)
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
			if game.Publisher != "Square Enix" {
				t.Errorf("expected publisher 'Square Enix', got '%s'", game.Publisher)
			}
		}
	})

	t.Run("ジャンルフィルタが動作すること", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/games?genre=RPG", nil)
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
			if game.Genre == nil || *game.Genre != "RPG" {
				t.Errorf("expected genre 'RPG', got %v", game.Genre)
			}
		}
	})

	t.Run("シリーズフィルタが動作すること", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/games?series=Final+Fantasy", nil)
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
			if game.Series == nil || *game.Series != "Final Fantasy" {
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
		req := httptest.NewRequest(http.MethodGet, "/games?platform=Nintendo+Switch&publisher=Nintendo", nil)
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
			if game.Platform != "Nintendo Switch" {
				t.Errorf("expected platform 'Nintendo Switch', got '%s'", game.Platform)
			}
			if game.Publisher != "Nintendo" {
				t.Errorf("expected publisher 'Nintendo', got '%s'", game.Publisher)
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
