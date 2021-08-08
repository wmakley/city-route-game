package city

import (
	"city-route-game/domain"
	"fmt"
	"os"
	"testing"

	"github.com/assertgo/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestMain(m *testing.M) {
	dbPath := "file::memory:?cache=shared"

	dbConn, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger:                   logger.Default.LogMode(logger.Error),
		DisableNestedTransaction: true,
	})
	if err != nil {
		panic("Error connecting to test db: " + err.Error())
	}

	if err = dbConn.AutoMigrate(domain.Models()...); err != nil {
		panic("Error migrating test db: " + err.Error())
	}

	Init(dbConn)

	os.Exit(m.Run())
}

func TestFindAllByBoardID(t *testing.T) {
	assert := assert.New(t)

	board := createTestBoard()
	city1 := createTestCity(board.ID)
	city2 := createTestCity(board.ID)

	results, err := FindAllByBoardID(board.ID)
	assert.That(err).IsNil()
	assert.ThatInt(len(results)).IsEqualTo(2)
	assert.That(results[0].ID).IsEqualTo(city1.ID)
	assert.That(results[1].ID).IsEqualTo(city2.ID)
}

func TestCreate(t *testing.T) {
	assert := assert.New(t)

	board := createTestBoard()

	form := Form{
		Name: "New City Name",
		Position: domain.Position{
			X: 123,
			Y: 432,
		},
	}

	city, err := Create(board.ID, &form)
	assert.That(err).IsNil()

	var updatedCity domain.City
	err = db.First(&updatedCity, city.ID).Error
	assert.That(err).IsNil()

	if updatedCity.Name != "New City Name" {
		t.Error("City Name was not updated")
	}
	if updatedCity.Position.X != 123 {
		t.Error("City Position X was not updated")
	}
	if updatedCity.Position.Y != 432 {
		t.Error("City Position Y was not updated")
	}
}

func TestAddSpace(t *testing.T) {
	assert := assert.New(t)
	board := createTestBoard()
	city := createTestCity(board.ID)

	form := AddCitySpaceForm{
		CityID:            city.ID,
		SpaceType:         domain.MerchantID,
		RequiredPrivilege: 2,
	}

	space, err := AddSpace(&form)
	if err != nil {
		t.Fatalf("Errors: %+v", form.Errors)
	}
	if space == nil {
		t.Fatal("space is nil")
	}
	assert.That(space.SpaceType).IsEqualTo(domain.MerchantID)
	assert.ThatInt(space.RequiredPrivilege).IsEqualTo(2)
	assert.ThatInt(space.Order).IsEqualTo(1)
}

type TestData struct {
	EmptyBoard            domain.Board
	BoardWithCities       domain.Board
	BoardWithCitiesCities []domain.City
}

func insertTestData() *TestData {
	emptyBoard := *createTestBoard()
	boardWithCities := *createTestBoard()

	cities := make([]domain.City, 0, 2)
	for i := 0; i < 2; i++ {
		cities = append(cities, *createTestCity(boardWithCities.ID))
	}

	return &TestData{
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
