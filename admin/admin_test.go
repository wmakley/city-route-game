package admin

import (
	"city-route-game/domain"
	"city-route-game/httpassert"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestAdmin(t *testing.T) {
	var err error

	dbPath := "../data/admin-test.sqlite"

	err = os.Remove(dbPath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal("Error deleting prior test database: ", err.Error())
	}

	var db *gorm.DB
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		t.Fatal("Error connecting to database: ", err.Error())
	}

	err = db.AutoMigrate(domain.Models()...)
	if err != nil {
		t.Fatal("Error migrating database: ", err.Error())
	}

	Init(db, "../templates")
	testData := insertTestData(t)

	router := NewAdminRouter()

	t.Run("getBoardByIdAsJson_includesCities", func(t *testing.T) {
		board := testData.BoardWithCities

		req := httptest.NewRequest("GET", fmt.Sprintf("/boards/%d", board.ID), nil)
		req.Header.Set("Accept", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		httpassert.Success(t, w)
		httpassert.JsonObject(t, w)

		responseJson := domain.Board{}
		if err := json.NewDecoder(w.Body).Decode(&responseJson); err != nil {
			t.Fatal(err)
		}

		if responseJson.ID != board.ID {
			t.Error("response ID does not match board ID")
		}

		if len(responseJson.Cities) != len(board.Cities) {
			t.Error("Cities were not returned")
		}
	})

	t.Run("listCitiesByBoardId_boardNotFound", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/boards/9999/cities/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		httpassert.NotFound(t, w)
	})

	t.Run("listCitiesByBoardId", func(t *testing.T) {
		board := createTestBoard(t)

		url := fmt.Sprintf("/boards/%d/cities/", board.ID)
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		httpassert.Success(t, w)
		httpassert.JsonArray(t, w)
	})
}

type TestData struct {
	EmptyBoard      domain.Board
	BoardWithCities domain.Board
}

func insertTestData(t *testing.T) TestData {
	emptyBoard := *createTestBoard(t)
	boardWithCities := *createTestBoard(t)

	for i := 0; i < 2; i++ {
		createTestCity(t, boardWithCities.ID)
	}

	return TestData{
		EmptyBoard:      emptyBoard,
		BoardWithCities: boardWithCities,
	}
}

var testBoardCounter = 0

func createTestBoard(t *testing.T) *domain.Board {
	board := domain.Board{
		Name: fmt.Sprintf("Test Board %d", testBoardCounter),
	}
	testBoardCounter++
	if err := db.Save(&board).Error; err != nil {
		t.Fatal(err)
	}
	return &board
}

func createTestCity(t *testing.T, boardID uint) *domain.City {
	city := domain.City{
		BoardID: boardID,
		Name:    "Test City",
	}

	if err := db.Save(&city).Error; err != nil {
		t.Fatal(err)
	}

	return &city
}
