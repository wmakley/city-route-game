package admin

import (
	"city-route-game/domain"
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

	err = db.AutoMigrate(&domain.Game{}, &domain.Board{}, &domain.Player{}, &domain.PlayerBoard{}, &domain.PlayerBonusToken{}, &domain.BonusToken{}, &domain.RouteBonusToken{}, &domain.City{}, &domain.CitySlot{}, &domain.Route{}, &domain.RouteSlot{})
	if err != nil {
		t.Fatal("Error migrating database: ", err.Error())
	}

	Init(db)
	router := NewAdminRouter()

	t.Run("listCitiesByBoardId", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:8080/boards/1/cities/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Error("Response code is not 200, was:", w.Code)
		}
	})
}
