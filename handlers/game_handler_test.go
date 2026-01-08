package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"takase-game-list/models"

	"takase-game-list/db"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// setupTestDB テスト用のインメモリデータベースをセットアップ
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	// 一時ファイルを使用してデータベースを作成
	testDB, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	if err := db.Migrate(testDB); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}
	return testDB
}

func TestCreateGame_Unit(t *testing.T) {
	db := setupTestDB(t)
	gin.SetMode(gin.TestMode)

	t.Run("正常なゲーム登録", func(t *testing.T) {
		router := gin.New()
		router.POST("/games", func(c *gin.Context) {
			CreateGame(c, db)
		})

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

		assert.Equal(t, http.StatusCreated, w.Code)

		var game models.Game
		err := json.Unmarshal(w.Body.Bytes(), &game)
		assert.NoError(t, err)
		assert.NotZero(t, game.ID)
		assert.Equal(t, "Test Game", game.Title)
	})

	t.Run("バリデーションエラー", func(t *testing.T) {
		router := gin.New()
		router.POST("/games", func(c *gin.Context) {
			CreateGame(c, db)
		})

		reqBody := map[string]interface{}{
			"title": "", // タイトルが空
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetGames_Unit(t *testing.T) {
	db := setupTestDB(t)
	gin.SetMode(gin.TestMode)

	// テストデータの準備
	games := []models.Game{
		{Title: "Game 1", ReleaseYear: 2024, Publisher: "Publisher 1", Platform: "PC"},
		{Title: "Game 2", ReleaseYear: 2023, Publisher: "Publisher 2", Platform: "Nintendo Switch"},
	}
	for i := range games {
		db.Create(&games[i])
	}

	router := gin.New()
	router.GET("/games", func(c *gin.Context) {
		GetGames(c, db)
	})

	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []models.Game
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestUpdateGame_Unit(t *testing.T) {
	db := setupTestDB(t)
	gin.SetMode(gin.TestMode)

	// テストデータの準備
	game := models.Game{
		Title:       "Original Game",
		ReleaseYear: 2024,
		Publisher:   "Original Publisher",
		Platform:    "PC",
	}
	db.Create(&game)

	router := gin.New()
	router.PUT("/games/:id", func(c *gin.Context) {
		UpdateGame(c, db)
	})

	t.Run("部分更新", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title": "Updated Game",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/games/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var updatedGame models.Game
		err := json.Unmarshal(w.Body.Bytes(), &updatedGame)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Game", updatedGame.Title)
		assert.Equal(t, "Original Publisher", updatedGame.Publisher)
	})

	t.Run("存在しないID", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title": "Updated Game",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/games/99999", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDeleteGame_Unit(t *testing.T) {
	db := setupTestDB(t)
	gin.SetMode(gin.TestMode)

	// テストデータの準備
	game := models.Game{
		Title:       "Game to Delete",
		ReleaseYear: 2024,
		Publisher:   "Publisher",
		Platform:    "PC",
	}
	db.Create(&game)

	router := gin.New()
	router.DELETE("/games/:id", func(c *gin.Context) {
		DeleteGame(c, db)
	})

	t.Run("削除成功", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/games/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// 削除されたことを確認
		var deletedGame models.Game
		result := db.First(&deletedGame, 1)
		assert.Error(t, result.Error)
		assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
	})

	t.Run("存在しないID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/games/99999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
