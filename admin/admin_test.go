package admin

import (
	"bytes"
	"city-route-game/domain"
	"city-route-game/httpassert"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	router   *mux.Router
	testData TestData
)

func TestMain(m *testing.M) {
	var err error

	dbPath := "../data/admin-test.sqlite"

	err = os.Remove(dbPath)
	if err != nil && !os.IsNotExist(err) {
		panic("Error deleting prior test database: " + err.Error())
	}

	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}

	err = db.AutoMigrate(domain.Models()...)
	if err != nil {
		panic("Error migrating database: " + err.Error())
	}

	Init(db, "../templates")
	testData = insertTestData()
	router = NewAdminRouter(false)

	os.Exit(m.Run())
}

func TestListBoards(t *testing.T) {
	req := httptest.NewRequest("GET", "/boards/", nil)
	req.Header.Set("Accept", "text/html")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	httpassert.Success(t, w)
	httpassert.HtmlContentType(t, w)
}

func TestNewBoard(t *testing.T) {
	req := httptest.NewRequest("GET", "/boards/new", nil)
	req.Header.Set("Accept", "text/html")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	httpassert.Success(t, w)
	httpassert.HtmlContentType(t, w)
}

func TestCreateBoard(t *testing.T) {
	formData := url.Values{}
	formData.Add("Name", fmt.Sprintf("Test Board %d", testBoardCounter))
	testBoardCounter++

	encodedFormData := formData.Encode()
	// t.Log("Encoded Form Data:", encodedFormData)

	req := httptest.NewRequest("POST", "/boards/", strings.NewReader(encodedFormData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("Accept", "text/html, text/javascript")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	httpassert.Success(t, w)
	httpassert.JavascriptContentType(t, w)
}

func TestGetBoardById_as_html(t *testing.T) {
	board := testData.EmptyBoard

	req := httptest.NewRequest("GET", fmt.Sprintf("/boards/%d", board.ID), nil)
	req.Header.Set("Accept", "text/html")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	httpassert.Success(t, w)
	httpassert.HtmlContentType(t, w)
}

func TestGetBoardById_asJson(t *testing.T) {
	board := testData.EmptyBoard

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

	if responseJson.Width != board.Width {
		t.Error("response Width does not match board")
	}

	if responseJson.Height != board.Height {
		t.Error("response Height does not match board")
	}
}

func TestEditBoard(t *testing.T) {
	board := testData.BoardWithCities
	req := httptest.NewRequest("GET", fmt.Sprintf("/boards/%d/edit", board.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Log("Body:", w.Body.String())

	httpassert.Success(t, w)
	httpassert.HtmlContentType(t, w)
}

func Test_update_board_name_via_web_form(t *testing.T) {
	board := createTestBoard()

	postData := url.Values{}
	postData.Set("_method", "PATCH")
	postData.Set("ID", fmt.Sprint(board.ID))
	postData.Set("Name", "New Name")

	req := httptest.NewRequest("POST", fmt.Sprintf("/boards/%d", board.ID), strings.NewReader(postData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("Accept", "text/javascript")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Log("Body:", w.Body)

	httpassert.Success(t, w)
	httpassert.JavascriptContentType(t, w)

	if err := db.Find(&board, board.ID).Error; err != nil {
		t.Fatalf("%+v", err)
	}

	if board.Name != "New Name" {
		t.Errorf("Board name was not updated (was '%s')", board.Name)
	}
}

func Test_update_board_dimensions_via_json(t *testing.T) {
	board := createTestBoard()

	newWidth, newHeight := 1234, 343

	payload := make(map[string]interface{})
	payload["id"] = board.ID
	payload["width"] = newWidth
	payload["height"] = newHeight

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("PATCH", fmt.Sprintf("/boards/%d", board.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if !httpassert.Success(t, w) {
		t.Log("Body: ", w.Body)
	}
	httpassert.JsonContentType(t, w)

	if err = db.Find(&board, board.ID).Error; err != nil {
		panic(err)
	}

	if board.Width != newWidth {
		t.Errorf("Board with was not updated (was %d)", board.Width)
	}
	if board.Height != newHeight {
		t.Errorf("Board height was was not updated (was %d)", board.Height)
	}
}

func TestListCitiesByBoardId_boardNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/boards/9999/cities/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	httpassert.NotFound(t, w)
}

func TestListCitiesByBoardId(t *testing.T) {
	board := testData.BoardWithCities

	url := fmt.Sprintf("/boards/%d/cities/", board.ID)
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	httpassert.Success(t, w)
	httpassert.JsonArray(t, w)

	responseJson := make([]domain.City, 0, len(testData.BoardWithCitiesCities))
	if err := json.NewDecoder(w.Body).Decode(&responseJson); err != nil {
		panic(err)
	}

	if len(responseJson) != len(testData.BoardWithCitiesCities) {
		t.Error("number of cities in json doesn't match number of cities in board")
	}
}

type TestData struct {
	EmptyBoard            domain.Board
	BoardWithCities       domain.Board
	BoardWithCitiesCities []domain.City
}

func insertTestData() TestData {
	emptyBoard := *createTestBoard()
	boardWithCities := *createTestBoard()

	cities := make([]domain.City, 0, 2)
	for i := 0; i < 2; i++ {
		cities = append(cities, *createTestCity(boardWithCities.ID))
	}

	return TestData{
		EmptyBoard:            emptyBoard,
		BoardWithCities:       boardWithCities,
		BoardWithCitiesCities: cities,
	}
}

var testBoardCounter = 0

func createTestBoard() *domain.Board {
	board := domain.Board{
		Name:   fmt.Sprintf("Test Board %d", testBoardCounter),
		Width:  10,
		Height: 20,
	}
	testBoardCounter++
	if err := db.Save(&board).Error; err != nil {
		panic(err)
	}
	return &board
}

func createTestCity(boardID uint) *domain.City {
	city := domain.City{
		BoardID: boardID,
		Name:    "Test City",
	}

	if err := db.Save(&city).Error; err != nil {
		panic(err)
	}

	return &city
}
