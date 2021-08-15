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
	db       *gorm.DB // original gorm connection
	rollback = errors.New("rollback")
)

func TestMain(m *testing.M) {
	dbPath := "file::memory:?cache=shared"
	var err error

	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
		DisableNestedTransaction: true,
	})
	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}

	if err := db.AutoMigrate(domain.Models()...); err != nil {
		panic("Error migrating database: " + err.Error())
	}

	os.Exit(m.Run())
}

// Create a transaction within which the gormBoardRepository interface may be tested using gorm itself
func TempTransaction(callback func (domain.BoardRepository, *gorm.DB)) {
	err := db.Transaction(func(tx *gorm.DB) error {
		repo := NewGormBoardRepository(tx)
		callback(repo, tx)
		return rollback
	})
	if err != nil && !errors.Is(err, rollback) {
		// should never happen
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
			err := p.CreateBoard(&board)
			if err != nil {
				t.Fatalf("CreateBoard failed: %+v", err)
			}
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

		err := p.CreateBoard(board)
		if err != nil {
			t.Errorf("CreateBoard returned error: %+v", err)
		}

		board, err = p.GetBoardByID(board.ID)
		if err != nil {
			t.Errorf("GetBoardByID %v returned error: %+v", board.ID, err)
		}

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
		err := p.CreateBoard(board)
		if err != nil {
			t.Fatalf("CreateBoard returned error: %+v", err)
		}

		originalID := board.ID

		err = p.UpdateBoard(originalID, func (board *domain.Board) (*domain.Board, error) {
			board.Name = "New Name"
			board.Width = 123
			board.Height = 321
			return board, nil
		})
		if err != nil {
			t.Fatalf("UpdateBoard returned error: %+v", err)
		}

		board, err = p.GetBoardByID(board.ID)
		if err != nil {
			t.Fatalf("GetBoardByID returned error: %+v", err)
		}

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
		err := p.CreateBoard(&board)
		assert.That(err).IsNil()

		city := domain.City{
			BoardID: board.ID,
			Name: "Test City",
		}
		err = p.CreateCity(&city)
		assert.That(err).IsNil()

		space := domain.CitySpace{
			CityID:    city.ID,
			Order:     1,
			SpaceType: domain.TraderID,
		}
		err = p.CreateCitySpace(&space)
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
		board := createTestBoard()
		createTestCityWithSpaces(board.ID)

		results, err := p.ListCitiesByBoardID(board.ID)
		assert.That(err).IsNil()
		assert.ThatInt(len(results)).IsEqualTo(1)
		assert.ThatInt(len(results[0].CitySpaces)).IsGreaterThan(0)
	})
}

func TestCreateCity(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func(p domain.BoardRepository, tx *gorm.DB) {
		board := createTestBoard()

		city := domain.City{
			BoardID: board.ID,
			Name: "New City Name",
			Position: domain.Position{
				X: 123,
				Y: 432,
			},
		}

		err := p.CreateCity(&city)
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

func TestUpdateCity(t *testing.T) {
	TempTransaction(func(p domain.BoardRepository, tx *gorm.DB) {
		board := createTestBoard()
		city := createTestCity(board.ID)

		err := p.UpdateCity(city, func(city *domain.City) (*domain.City, error) {
			city.Name = "New City Name"
			city.Position.X = 123
			city.Position.Y = 432
			return city, nil
		})
		if err != nil {
			t.Fatalf("UpdateCity returned error: %+v", err)
		}

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
		err := p.CreateCitySpace(&space)
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
		err := p.CreateCitySpace(&space)
		assert.That(err).IsNil()

		err = p.DeleteCitySpaceByID(space.ID)
		assert.That(err).IsNil()

		city, err = p.GetCityByID(city.ID)
		assert.That(err).IsNil()

		assert.ThatInt(len(city.CitySpaces)).IsEqualTo(0)
	})
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

	if err := db.Save(city.CitySpaces).Error; err != nil {
		panic(err)
	}

	return city
}
