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

// setupTestDBForMaster マスタデータハンドラーテスト用のデータベースをセットアップ
func setupTestDBForMaster(t *testing.T) *gorm.DB {
	t.Helper()
	testDB, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	if err := db.Migrate(testDB); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}
	return testDB
}

func TestGetPublishers_Unit(t *testing.T) {
	db := setupTestDBForMaster(t)
	gin.SetMode(gin.TestMode)

	// テストデータの準備
	publisher1 := models.Publisher{Name: "Publisher 1"}
	publisher2 := models.Publisher{Name: "Publisher 2"}
	db.Create(&publisher1)
	db.Create(&publisher2)

	router := gin.New()
	router.GET("/publishers", func(c *gin.Context) {
		GetPublishers(c, db)
	})

	req := httptest.NewRequest(http.MethodGet, "/publishers", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []models.Publisher
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCreatePublisher_Unit(t *testing.T) {
	db := setupTestDBForMaster(t)
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/publishers", func(c *gin.Context) {
		CreatePublisher(c, db)
	})

	t.Run("正常なPublisher作成", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "New Publisher",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/publishers", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var publisher models.Publisher
		err := json.Unmarshal(w.Body.Bytes(), &publisher)
		assert.NoError(t, err)
		assert.NotZero(t, publisher.ID)
		assert.Equal(t, "New Publisher", publisher.Name)
	})

	t.Run("Nameが空の場合", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/publishers", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("重複したNameの場合", func(t *testing.T) {
		// 既存のPublisherを作成
		existingPublisher := models.Publisher{Name: "Existing Publisher"}
		db.Create(&existingPublisher)

		reqBody := map[string]interface{}{
			"name": "Existing Publisher",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/publishers", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestGetPlatforms_Unit(t *testing.T) {
	db := setupTestDBForMaster(t)
	gin.SetMode(gin.TestMode)

	// テストデータの準備
	platform1 := models.Platform{Name: "PC"}
	platform2 := models.Platform{Name: "Nintendo Switch"}
	db.Create(&platform1)
	db.Create(&platform2)

	router := gin.New()
	router.GET("/platforms", func(c *gin.Context) {
		GetPlatforms(c, db)
	})

	req := httptest.NewRequest(http.MethodGet, "/platforms", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []models.Platform
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCreatePlatform_Unit(t *testing.T) {
	db := setupTestDBForMaster(t)
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/platforms", func(c *gin.Context) {
		CreatePlatform(c, db)
	})

	t.Run("正常なPlatform作成", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "New Platform",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/platforms", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var platform models.Platform
		err := json.Unmarshal(w.Body.Bytes(), &platform)
		assert.NoError(t, err)
		assert.NotZero(t, platform.ID)
		assert.Equal(t, "New Platform", platform.Name)
	})

	t.Run("Nameが空の場合", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/platforms", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("重複したNameの場合", func(t *testing.T) {
		// 既存のPlatformを作成
		existingPlatform := models.Platform{Name: "Existing Platform"}
		db.Create(&existingPlatform)

		reqBody := map[string]interface{}{
			"name": "Existing Platform",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/platforms", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestGetSeries_Unit(t *testing.T) {
	db := setupTestDBForMaster(t)
	gin.SetMode(gin.TestMode)

	// テストデータの準備
	series1 := models.Series{Name: "Series 1"}
	series2 := models.Series{Name: "Series 2"}
	db.Create(&series1)
	db.Create(&series2)

	router := gin.New()
	router.GET("/series", func(c *gin.Context) {
		GetSeries(c, db)
	})

	req := httptest.NewRequest(http.MethodGet, "/series", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []models.Series
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCreateSeries_Unit(t *testing.T) {
	db := setupTestDBForMaster(t)
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/series", func(c *gin.Context) {
		CreateSeries(c, db)
	})

	t.Run("正常なSeries作成", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "New Series",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/series", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var series models.Series
		err := json.Unmarshal(w.Body.Bytes(), &series)
		assert.NoError(t, err)
		assert.NotZero(t, series.ID)
		assert.Equal(t, "New Series", series.Name)
	})

	t.Run("Nameが空の場合", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/series", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("重複したNameの場合", func(t *testing.T) {
		// 既存のSeriesを作成
		existingSeries := models.Series{Name: "Existing Series"}
		db.Create(&existingSeries)

		reqBody := map[string]interface{}{
			"name": "Existing Series",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/series", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestGetGenres_Unit(t *testing.T) {
	db := setupTestDBForMaster(t)
	gin.SetMode(gin.TestMode)

	// テストデータの準備
	genre1 := models.Genre{Name: "RPG"}
	genre2 := models.Genre{Name: "Action"}
	db.Create(&genre1)
	db.Create(&genre2)

	router := gin.New()
	router.GET("/genres", func(c *gin.Context) {
		GetGenres(c, db)
	})

	req := httptest.NewRequest(http.MethodGet, "/genres", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []models.Genre
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCreateGenre_Unit(t *testing.T) {
	db := setupTestDBForMaster(t)
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/genres", func(c *gin.Context) {
		CreateGenre(c, db)
	})

	t.Run("正常なGenre作成", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "New Genre",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/genres", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var genre models.Genre
		err := json.Unmarshal(w.Body.Bytes(), &genre)
		assert.NoError(t, err)
		assert.NotZero(t, genre.ID)
		assert.Equal(t, "New Genre", genre.Name)
	})

	t.Run("Nameが空の場合", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/genres", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("重複したNameの場合", func(t *testing.T) {
		// 既存のGenreを作成
		existingGenre := models.Genre{Name: "Existing Genre"}
		db.Create(&existingGenre)

		reqBody := map[string]interface{}{
			"name": "Existing Genre",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/genres", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}
