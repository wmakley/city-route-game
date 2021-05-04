package admin

import (
	"city-route-game/domain"
	"fmt"
	"net/http/httptest"
	"os"
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

	Init(db)
	router := NewAdminRouter()

	testBoardCounter := 0

	t.Run("listCitiesByBoardId_boardNotFound", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:8080/boards/9999/cities/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Error("Expected response code to be 404 because board 9999 doesn't exist, but was:", w.Code, "body:", w.Body)
		}
	})

	t.Run("listCitiesByBoardId", func(t *testing.T) {
		board := domain.Board{
			Name: fmt.Sprintf("Test Board %d", testBoardCounter),
		}
		testBoardCounter++
		if err := db.Save(&board).Error; err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("http://localhost:8080/boards/%d/cities/", board.ID)
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Error("Response code is not 200, was:", w.Code, "body:", w.Body)
		}
	})
}
