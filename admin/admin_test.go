package admin

import (
	"bytes"
	"city-route-game/httpassert"
	"city-route-game/internal/app"
	"city-route-game/internal/gorm_board_crud_repository"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/schema"
	"net/http"
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
	repo               app.BoardCrudRepository
	boardEditorService app.BoardEditorService
)

func TestMain(m *testing.M) {
	dbPath := "../data/admin-test.sqlite"

	err := os.Remove(dbPath)
	if err != nil && !os.IsNotExist(err) {
		panic("Error deleting prior test gorm_board_crud_repository: " + err.Error())
	}

	dbConn, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Error connecting to gorm_board_crud_repository: " + err.Error())
	}

	err = dbConn.AutoMigrate(gorm_board_crud_repository.Models()...)
	if err != nil {
		panic("Error migrating gorm_board_crud_repository: " + err.Error())
	}

	repo = gorm_board_crud_repository.NewGormBoardCrudRepository(dbConn)
	boardEditorService = app.NewBoardEditorService(repo)

	controllerConfig := ControllerConfig{
		FormDecoder: schema.NewDecoder(),
		TemplateRoot: "../templates",
		AssetHost: "",
	}

	boardController := NewBoardController(controllerConfig, boardEditorService)
	cityController := NewCityController(controllerConfig, boardEditorService)

	testData = insertTestData(context.Background())
	router = NewAdminRouter(&boardController, &cityController, []string{}, false)

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
	formData.Add("name", fmt.Sprintf("Test Board %d", testBoardCounter))
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

	responseJson := app.Board{}
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

	//t.Log("Body:", w.Body.String())

	httpassert.Success(t, w)
	httpassert.HtmlContentType(t, w)
}

func Test_update_board_name_via_web_form(t *testing.T) {
	ctx := context.Background()
	board := createTestBoard(ctx)

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

	//updatedBoard, err := repo.GetBoardByID(board.ID)
	//if err != nil {
	//	t.Fatalf("%+v", err)
	//}
	//
	//if updatedBoard.Name != "New Name" {
	//	t.Errorf("Board name was not updated (was '%s')", board.Name)
	//}
}

func Test_update_board_dimensions_via_json(t *testing.T) {
	ctx := context.Background()
	board := createTestBoard(ctx)

	newWidth, newHeight := 1234, 343

	payload := make(map[string]interface{})
	payload["name"] = board.Name
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

	updatedBoard, err := repo.GetBoardByID(ctx, board.ID)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if updatedBoard.Width != newWidth {
		t.Errorf("Board with was not updated (was %d)", board.Width)
	}
	if updatedBoard.Height != newHeight {
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

	responseJson := make([]app.City, 0, len(testData.BoardWithCitiesCities))
	if err := json.NewDecoder(w.Body).Decode(&responseJson); err != nil {
		panic(err)
	}

	if len(responseJson) != len(testData.BoardWithCitiesCities) {
		t.Error("number of cities in json doesn't match number of cities in board")
	}
}

func TestCreateCity(t *testing.T) {
	ctx := context.Background()
	board := createTestBoard(ctx)
	url := fmt.Sprintf("/boards/%d/cities/", board.ID)
	city := app.CityForm{
		Name: "Test City",
		Position: app.Position{
			X: 10,
			Y: 20,
		},
	}

	body, err := json.Marshal(&city)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if !httpassert.Success(t, w) {
		t.Log("Body:", w.Body)
	}
	httpassert.JsonContentType(t, w)
}

func TestUpdateCity(t *testing.T) {
	ctx := context.Background()
	board := createTestBoard(ctx)
	city := createTestCity(ctx, board.ID)
	url := fmt.Sprintf("/boards/%d/cities/%d", board.ID, city.ID)

	newName := "New City Name"
	newX := 123
	newY := 432
	form := app.CityForm{
		Name: newName,
		Position: app.Position{
			X: newX,
			Y: newY,
		},
	}

	body, err := json.Marshal(&form)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest("PUT", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if !httpassert.Success(t, w) {
		t.Log("Body:", w.Body)
	}
	httpassert.JsonContentType(t, w)

	var updatedCity app.City
	if err = json.NewDecoder(w.Body).Decode(&updatedCity); err != nil {
		panic(err)
	}

	if updatedCity.Name != newName {
		t.Error("City Name was not updated")
	}
	if updatedCity.Position.X != newX {
		t.Error("City Position X was not updated")
	}
	if updatedCity.Position.Y != newY {
		t.Error("City Position Y was not updated")
	}
}

func TestDeleteBoard(t *testing.T) {
	ctx := context.Background()
	board := createTestBoard(ctx)
	city := createTestCity(ctx, board.ID)
	space := app.CitySpace{
		CityID:    city.ID,
		Order:     1,
		SpaceType: app.TraderID,
	}
	var err error
	if err = repo.CreateCitySpace(ctx, &space); err != nil {
		t.Fatalf("Error saving test space: %+v", err)
	}
	if city, err = repo.GetCityByID(ctx, city.ID); err != nil {
		t.Fatalf("Error reloading city: %+v", err)
	}
	if len(city.CitySpaces) != 1 {
		t.Error("CitySpaces HasMany relationship is not loading")
	}

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/boards/%d", board.ID), nil)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if !httpassert.Success(t, w) {
		t.Log("Body:", w.Body)
	}
	httpassert.JavascriptContentType(t, w)
}

func TestDeleteCity(t *testing.T) {
	ctx := context.Background()
	board := createTestBoard(ctx)
	city := createTestCity(ctx, board.ID)
	space := app.CitySpace{
		CityID:    city.ID,
		Order:     1,
		SpaceType: app.TraderID,
	}
	if err := repo.CreateCitySpace(ctx, &space); err != nil {
		t.Fatalf("Error saving test space: %+v", err)
	}
	var err error
	if city, err = repo.GetCityByID(ctx, city.ID); err != nil {
		t.Fatalf("Error reloading city: %+v", err)
	}
	if len(city.CitySpaces) != 1 {
		t.Error("CitySpaces HasMany relationship is not loading")
	}

	url := fmt.Sprintf("/boards/%d/cities/%d", board.ID, city.ID)
	req := httptest.NewRequest("DELETE", url, nil)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Response code is not 204 (is %d)", w.Code)
		t.Log("Body:", w.Body)
	}

	var cities []app.City
	if cities, err = repo.ListCitiesByBoardID(ctx, city.ID); err != nil {
		panic(err)
	}

	//var spaces []app.CitySpace
	//if err = repo.Find(&spaces, "city_id = ?", city.ID).Error; err != nil {
	//	panic(err)
	//}
	//if len(spaces) != 0 {
	//	t.Error("City spaces were not deleted")
	//}

	if len(cities) != 0 {
		t.Error("City was not deleted")
	}
}

type TestData struct {
	EmptyBoard            app.Board
	BoardWithCities       app.Board
	BoardWithCitiesCities []app.City
}

func insertTestData(ctx context.Context) TestData {
	emptyBoard := *createTestBoard(ctx)
	boardWithCities := *createTestBoard(ctx)

	cities := make([]app.City, 0, 2)
	for i := 0; i < 2; i++ {
		cities = append(cities, *createTestCity(ctx, boardWithCities.ID))
	}

	return TestData{
		EmptyBoard:            emptyBoard,
		BoardWithCities:       boardWithCities,
		BoardWithCitiesCities: cities,
	}
}

var testBoardCounter = 0

func createTestBoard(ctx context.Context) *app.Board {
	form := app.NewCreateBoardForm()
	form.Name = fmt.Sprintf("Test Board %d", testBoardCounter)
	testBoardCounter++
	board, err := boardEditorService.CreateBoard(ctx, &form)
	if err != nil {
		panic(err)
	}
	return board
}

func createTestCity(ctx context.Context, boardID app.ID) *app.City {
	city := app.City{
		BoardID: boardID,
		Name:    "Test City",
	}

	err := repo.CreateCity(ctx, &city)
	if err != nil {
		panic(err)
	}

	return &city
}
