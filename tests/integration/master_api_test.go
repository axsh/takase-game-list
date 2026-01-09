package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"takase-game-list/handlers"
	"takase-game-list/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// setupTestServerForMaster マスタデータAPIテスト用のHTTPサーバーをセットアップ
func setupTestServerForMaster(t *testing.T, database *gorm.DB) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// マスタデータAPIのルーティング
	router.GET("/publishers", func(c *gin.Context) {
		handlers.GetPublishers(c, database)
	})
	router.POST("/publishers", func(c *gin.Context) {
		handlers.CreatePublisher(c, database)
	})
	router.GET("/platforms", func(c *gin.Context) {
		handlers.GetPlatforms(c, database)
	})
	router.POST("/platforms", func(c *gin.Context) {
		handlers.CreatePlatform(c, database)
	})
	router.GET("/series", func(c *gin.Context) {
		handlers.GetSeries(c, database)
	})
	router.POST("/series", func(c *gin.Context) {
		handlers.CreateSeries(c, database)
	})
	router.GET("/genres", func(c *gin.Context) {
		handlers.GetGenres(c, database)
	})
	router.POST("/genres", func(c *gin.Context) {
		handlers.CreateGenre(c, database)
	})

	return router
}

// TestGetPublishers_Integration Publisher一覧取得の確認
func TestGetPublishers_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServerForMaster(t, database)

	// テストデータの準備
	publisher1 := models.Publisher{Name: "Publisher 1"}
	publisher2 := models.Publisher{Name: "Publisher 2"}
	database.Create(&publisher1)
	database.Create(&publisher2)

	req := httptest.NewRequest(http.MethodGet, "/publishers", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result []models.Publisher
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 publishers, got %d", len(result))
	}
}

// TestCreatePublisher_Integration Publisher作成の確認
func TestCreatePublisher_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServerForMaster(t, database)

	reqBody := map[string]interface{}{
		"name": "New Publisher",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/publishers", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var publisher models.Publisher
	if err := json.Unmarshal(w.Body.Bytes(), &publisher); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if publisher.ID == 0 {
		t.Error("publisher ID was not auto-generated")
	}
	if publisher.Name != "New Publisher" {
		t.Errorf("expected name 'New Publisher', got '%s'", publisher.Name)
	}
}

// TestCreatePublisherDuplicate_Integration Publisher重複作成のエラー確認
func TestCreatePublisherDuplicate_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServerForMaster(t, database)

	// 既存のPublisherを作成
	existingPublisher := models.Publisher{Name: "Existing Publisher"}
	database.Create(&existingPublisher)

	reqBody := map[string]interface{}{
		"name": "Existing Publisher",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/publishers", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, w.Code)
	}
}

// TestGetPlatforms_Integration Platform一覧取得の確認
func TestGetPlatforms_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServerForMaster(t, database)

	// テストデータの準備
	platform1 := models.Platform{Name: "PC"}
	platform2 := models.Platform{Name: "Nintendo Switch"}
	database.Create(&platform1)
	database.Create(&platform2)

	req := httptest.NewRequest(http.MethodGet, "/platforms", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result []models.Platform
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 platforms, got %d", len(result))
	}
}

// TestCreatePlatform_Integration Platform作成の確認
func TestCreatePlatform_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServerForMaster(t, database)

	reqBody := map[string]interface{}{
		"name": "New Platform",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/platforms", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var platform models.Platform
	if err := json.Unmarshal(w.Body.Bytes(), &platform); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if platform.ID == 0 {
		t.Error("platform ID was not auto-generated")
	}
	if platform.Name != "New Platform" {
		t.Errorf("expected name 'New Platform', got '%s'", platform.Name)
	}
}

// TestCreatePlatformDuplicate_Integration Platform重複作成のエラー確認
func TestCreatePlatformDuplicate_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServerForMaster(t, database)

	// 既存のPlatformを作成
	existingPlatform := models.Platform{Name: "Existing Platform"}
	database.Create(&existingPlatform)

	reqBody := map[string]interface{}{
		"name": "Existing Platform",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/platforms", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, w.Code)
	}
}

// TestGetSeries_Integration Series一覧取得の確認
func TestGetSeries_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServerForMaster(t, database)

	// テストデータの準備
	series1 := models.Series{Name: "Series 1"}
	series2 := models.Series{Name: "Series 2"}
	database.Create(&series1)
	database.Create(&series2)

	req := httptest.NewRequest(http.MethodGet, "/series", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result []models.Series
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 series, got %d", len(result))
	}
}

// TestCreateSeries_Integration Series作成の確認
func TestCreateSeries_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServerForMaster(t, database)

	reqBody := map[string]interface{}{
		"name": "New Series",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/series", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var series models.Series
	if err := json.Unmarshal(w.Body.Bytes(), &series); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if series.ID == 0 {
		t.Error("series ID was not auto-generated")
	}
	if series.Name != "New Series" {
		t.Errorf("expected name 'New Series', got '%s'", series.Name)
	}
}

// TestGetGenres_Integration Genre一覧取得の確認
func TestGetGenres_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServerForMaster(t, database)

	// テストデータの準備
	genre1 := models.Genre{Name: "RPG"}
	genre2 := models.Genre{Name: "Action"}
	database.Create(&genre1)
	database.Create(&genre2)

	req := httptest.NewRequest(http.MethodGet, "/genres", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result []models.Genre
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 genres, got %d", len(result))
	}
}

// TestCreateGenre_Integration Genre作成の確認
func TestCreateGenre_Integration(t *testing.T) {
	database, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, database, dbPath)

	router := setupTestServerForMaster(t, database)

	reqBody := map[string]interface{}{
		"name": "New Genre",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/genres", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var genre models.Genre
	if err := json.Unmarshal(w.Body.Bytes(), &genre); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if genre.ID == 0 {
		t.Error("genre ID was not auto-generated")
	}
	if genre.Name != "New Genre" {
		t.Errorf("expected name 'New Genre', got '%s'", genre.Name)
	}
}
