package admin

import (
	"city-route-game/domain"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAdmin(t *testing.T) {
	var err error

	dbPath := "../data/admin-test.sqlite"

	err = os.Remove(dbPath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal("Error deleting prior test database: ", err.Error())
	}

	var db *gorm.DB
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
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

		req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8080/boards/%d", board.ID), nil)
		req.Header.Set("Accept", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Error("Expected response to be 200, but was:", w.Code, "body:", w.Body)
		}

		if !strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") {
			t.Fatal("Content type is not JSON")
		}

		responseJson := domain.Board{}
		if err := json.NewDecoder(w.Body).Decode(&responseJson); err != nil {
			t.Fatal(err)
		}

		if responseJson.ID != board.ID {
			t.Error("response ID does not match board ID")
		}
	})

	t.Run("listCitiesByBoardId_boardNotFound", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:8080/boards/9999/cities/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Error("Expected response code to be 404 because board 9999 doesn't exist, but was:", w.Code, "body:", w.Body)
		}
	})

	t.Run("listCitiesByBoardId", func(t *testing.T) {
		board := createTestBoard(t)

		url := fmt.Sprintf("http://localhost:8080/boards/%d/cities/", board.ID)
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Error("Response code is not 200, was:", w.Code, "body:", w.Body)
		}
	})
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

type TestData struct {
	EmptyBoard      domain.Board
	BoardWithCities domain.Board
}

func insertTestData(t *testing.T) TestData {
	emptyBoard := *createTestBoard(t)
	boardWithCities := *createTestBoard(t)

	return TestData{
		EmptyBoard:      emptyBoard,
		BoardWithCities: boardWithCities,
	}
}
