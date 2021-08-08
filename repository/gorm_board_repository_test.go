package repository

import (
	"city-route-game/domain"
	"errors"
	"fmt"
	"github.com/assertgo/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"testing"
)

var (
	DB *gorm.DB // original gorm connection
	Rollback = errors.New("rollback")
)

func TestMain(m *testing.M) {
	dbPath := "file::memory:?cache=shared"
	var err error

	//err := os.Remove(dbPath)
	//if err != nil && !os.IsNotExist(err) {
	//	panic("Error deleting prior test repository: " + err.Error())
	//}

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		DisableNestedTransaction: true,
	})
	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}

	if err := DB.AutoMigrate(domain.Models()...); err != nil {
		panic("Error migrating database: " + err.Error())
	}

	os.Exit(m.Run())
}

// Create a transaction within which the GormRepository interface may be tested using gorm itself
func TempTransaction(callback func (domain.BoardRepository, *gorm.DB)) {
	err := DB.Transaction(func(tx *gorm.DB) error {
		repo := NewGormRepository(tx)
		callback(repo, tx)
		return Rollback
	})
	if err != nil && !errors.Is(err, Rollback) {
		panic(err)
	}
}

func TestListBoards(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func (p domain.BoardRepository, tx *gorm.DB) {
		boards := []domain.Board{
			{
				Name:   "Test Board 1",
				Width:  10,
				Height: 20,
			},
			{
				Name:   "Test Board 2",
				Width:  10,
				Height: 20,
			},
		}

		for _, board := range boards {
			err := p.SaveBoard(&board)
			assert.That(err).IsNil()
		}

		results, err := p.ListBoards()
		assert.That(err).IsNil()
		assert.ThatInt(len(results)).IsEqualTo(2)
		assert.ThatString(results[0].Name).IsEqualTo("Test Board 1")
		assert.ThatString(results[1].Name).IsEqualTo("Test Board 2")
	})
}

func TestCreateBoard(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func(p domain.BoardRepository, tx *gorm.DB) {
		board := &domain.Board{
			Name:   "My Awesome Board",
			Width:  200,
			Height: 300,
		}

		err := p.SaveBoard(board)

		board, err = p.GetBoardByID(board.ID)
		assert.That(err).IsNil()

		assert.That(err).IsNil()
		assert.That(board).IsNotNil()
		assert.ThatUint64(uint64(board.ID)).IsNonZero()
		assert.ThatString(board.Name).IsEqualTo("My Awesome Board")
		assert.ThatInt(board.Width).IsEqualTo(200)
		assert.ThatInt(board.Height).IsEqualTo(300)
	})
}

func TestUpdateBoard(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func(p domain.BoardRepository, tx *gorm.DB) {
		board := &domain.Board{
			Name:   "Original Name",
			Width:  200,
			Height: 300,
		}
		err := p.SaveBoard(board)
		assert.That(err).IsNil()

		originalID := board.ID

		board.Name = "New Name"
		board.Width = 123
		board.Height = 321

		err = p.SaveBoard(board)
		assert.That(err).IsNil()

		board, err = p.GetBoardByID(board.ID)
		assert.That(err).IsNil()

		assert.That(board.ID).IsEqualTo(originalID)
		assert.ThatString(board.Name).IsEqualTo("New Name")
		assert.ThatInt(board.Width).IsEqualTo(123)
		assert.ThatInt(board.Height).IsEqualTo(321)
	})
}

func TestDeleteBoardByIDDeletesNestedRecords(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func(p domain.BoardRepository, tx *gorm.DB) {
		board := domain.Board{
			Name: "Test Board",
		}
		err := p.SaveBoard(&board)
		assert.That(err).IsNil()

		city := domain.City{
			BoardID: board.ID,
			Name: "Test City",
		}
		err = p.SaveCity(&city)
		assert.That(err).IsNil()

		space := domain.CitySpace{
			CityID:    city.ID,
			Order:     1,
			SpaceType: domain.TraderID,
		}
		err = p.SaveCitySpace(&space)
		assert.That(err).IsNil()

		err = p.DeleteBoardByID(board.ID)
		assert.That(err).IsNil()

		var cities []domain.City
		err = tx.Find(&cities, "board_id = ?", board.ID).Error
		assert.That(err).IsNil()
		assert.ThatInt(len(cities)).IsEqualTo(0)

		var spaces []domain.CitySpace
		err = tx.Find(&spaces, "city_id = ?", city.ID).Error
		assert.That(err).IsNil()
		assert.ThatInt(len(spaces)).IsEqualTo(0)

		var boards []domain.Board
		err = tx.Find(&boards, board.ID).Error
		assert.That(err).IsNil()
		assert.ThatInt(len(boards)).IsEqualTo(0)
	})
}

func TestListCitiesByBoardId(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func (p domain.BoardRepository, tx *gorm.DB) {
		testData := insertTestData()
		createTestCityWithSpaces(testData.EmptyBoard.ID)
		board := testData.BoardWithCities

		results, err := p.ListCitiesByBoardID(board.ID)
		assert.That(err).IsNil()
		assert.ThatInt(len(results)).IsEqualTo(len(testData.BoardWithCitiesCities))
		assert.ThatInt(len(results[0].CitySpaces)).IsGreaterThan(0)
	})
}

func TestSaveCity(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func(p domain.BoardRepository, tx *gorm.DB) {
		board := createTestBoard()
		city := createTestCity(board.ID)

		city.Name = "New City Name"
		city.Position.X = 123
		city.Position.Y = 432

		err := p.SaveCity(city)
		assert.That(err).IsNil()

		var updatedCity domain.City
		if err := tx.Find(&updatedCity, city.ID).Error; err != nil {
			panic(err)
		}

		if updatedCity.Name != "New City Name" {
			t.Error("City Name was not updated")
		}
		if updatedCity.Position.X != 123 {
			t.Error("City Position X was not updated")
		}
		if updatedCity.Position.Y != 432 {
			t.Error("City Position Y was not updated")
		}
	})
}

func TestGetCitySpacesByCityID(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func (p domain.BoardRepository, tx *gorm.DB) {
		board := createTestBoard()
		city := createTestCityWithSpaces(board.ID)

		spaces, err := p.GetCitySpacesByCityID(city.ID)
		assert.That(err).IsNil()

		assert.ThatInt(len(spaces)).IsEqualTo(len(city.CitySpaces))
	})
}

func TestSaveCitySpace(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func (p domain.BoardRepository, tx *gorm.DB) {
		board := createTestBoard()
		city := createTestCity(board.ID)

		space := domain.CitySpace{
			CityID: city.ID,
			Order: 1,
			SpaceType: domain.MerchantID,
			RequiredPrivilege: 2,
		}
		err := p.SaveCitySpace(&space)
		assert.That(err).IsNil()

		city, err = p.GetCityByID(city.ID)
		assert.That(err).IsNil()

		space = city.CitySpaces[0]
		assert.That(space.SpaceType).IsEqualTo(domain.MerchantID)
		assert.ThatInt(space.RequiredPrivilege).IsEqualTo(2)
		assert.ThatInt(space.Order).IsEqualTo(1)
	})
}

func TestDeleteCitySpaceByIDSucceeds(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func (p domain.BoardRepository, tx *gorm.DB) {
		board := createTestBoard()
		city := createTestCity(board.ID)

		space := domain.CitySpace{
			CityID: city.ID,
			Order: 1,
			SpaceType: domain.MerchantID,
			RequiredPrivilege: 2,
		}
		err := p.SaveCitySpace(&space)
		assert.That(err).IsNil()

		p.DeleteCitySpaceByID(space.ID)

		city, err = p.GetCityByID(city.ID)
		assert.That(err).IsNil()

		space = city.CitySpaces[0]
		assert.That(space.SpaceType).IsEqualTo(domain.MerchantID)
		assert.ThatInt(space.RequiredPrivilege).IsEqualTo(2)
		assert.ThatInt(space.Order).IsEqualTo(1)
	})
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
	if err := DB.Save(&board).Error; err != nil {
		panic(err)
	}
	return &board
}

func createTestCity(boardID uint) *domain.City {
	city := domain.City{
		BoardID: boardID,
		Name:    "Test City",
	}

	if err := DB.Save(&city).Error; err != nil {
		panic(err)
	}

	return &city
}

func createTestCityWithSpaces(boardID uint) *domain.City {
	city := createTestCity(boardID)

	city.CitySpaces = []domain.CitySpace{
		{
			CityID:            city.ID,
			Order:             1,
			SpaceType:         domain.TraderID,
			RequiredPrivilege: 1,
		},
		{
			CityID:            city.ID,
			Order:             2,
			SpaceType:         domain.MerchantID,
			RequiredPrivilege: 2,
		},
		{
			CityID:            city.ID,
			Order:             3,
			SpaceType:         domain.TraderID,
			RequiredPrivilege: 3,
		},
	}

	if err := DB.Save(city.CitySpaces).Error; err != nil {
		panic(err)
	}

	return city
}
