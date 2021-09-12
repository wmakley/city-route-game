package gorm_board_crud_repository

import (
	"city-route-game/internal/app"
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

	if err := db.AutoMigrate(Models()...); err != nil {
		panic("Error migrating database: " + err.Error())
	}

	os.Exit(m.Run())
}

// Create a transaction within which the gormBoardRepository interface may be tested using gorm itself
func TempTransaction(callback func (app.BoardCrudRepository, *gorm.DB)) {
	err := db.Transaction(func(tx *gorm.DB) error {
		repo := NewGormBoardCrudRepository(tx)
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
	TempTransaction(func (p app.BoardCrudRepository, tx *gorm.DB) {
		boards := []app.Board{
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
		if err != nil {
			t.Fatalf("ListBoards returned error: %+v", err)
		}
		//t.Logf("Results: %+v", results)
		assert.ThatInt(len(results)).IsEqualTo(2)
		assert.ThatString(results[0].Name).IsEqualTo("Test Board 1")
		assert.ThatString(results[1].Name).IsEqualTo("Test Board 2")
	})
}

func TestCreateBoard(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func(p app.BoardCrudRepository, tx *gorm.DB) {
		board := &app.Board{
			Name:   "My Awesome Board",
			Width:  200,
			Height: 300,
		}

		err := p.CreateBoard(board)
		if err != nil {
			t.Fatalf("CreateBoard returned error: %+v", err)
		}

		boardID := board.ID
		board, err = p.GetBoardByID(boardID)
		if err != nil {
			t.Fatalf("GetBoardByID %v returned error: %+v", boardID, err)
		}

		assert.That(board).IsNotNil()
		assert.ThatUint64(uint64(board.ID)).IsNonZero()
		assert.ThatString(board.Name).IsEqualTo("My Awesome Board")
		assert.ThatInt(board.Width).IsEqualTo(200)
		assert.ThatInt(board.Height).IsEqualTo(300)
	})
}

func TestCreateBoardReturnsErrorOnDuplicateName(t *testing.T) {
	var beginCount int64
	err := db.Model(&Board{}).Count(&beginCount).Error
	if err != nil {
		panic(err)
	}

	assert := assert.New(t)
	TempTransaction(func(r app.BoardCrudRepository, tx *gorm.DB) {
		board1 := app.Board{
			Name:   "My Awesome Board",
			Width:  200,
			Height: 300,
		}

		err := r.CreateBoard(&board1)
		if err != nil {
			t.Fatalf("CreateBoard returned error: %+v", err)
		}

		board2 := app.Board{
			Name:   "My Awesome Board",
			Width:  200,
			Height: 300,
		}
		err = r.CreateBoard(&board2)
		if err != app.ErrNameTaken {
			t.Error("Expected ErrNameTaken to be returned for duplicate board name")
		}
	})

	var endCount int64
	err = db.Model(&Board{}).Count(&endCount).Error
	if err != nil {
		panic(err)
	}
	assert.That(endCount).IsEqualTo(beginCount)
}

func TestUpdateBoard(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func(p app.BoardCrudRepository, tx *gorm.DB) {
		board := &app.Board{
			Name:   "Original Name",
			Width:  200,
			Height: 300,
		}
		err := p.CreateBoard(board)
		if err != nil {
			t.Fatalf("CreateBoard returned error: %+v", err)
		}

		originalID := board.ID

		board, err = p.UpdateBoard(originalID, func (board *app.Board) (*app.Board, error) {
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
	TempTransaction(func(p app.BoardCrudRepository, tx *gorm.DB) {
		board := app.Board{
			Name: "Test Board",
		}
		err := p.CreateBoard(&board)
		if err != nil {
			t.Fatalf("CreateBoard returned error %+v", err)
		}

		city := app.City{
			BoardID: board.ID,
			Name: "Test City",
		}
		err = p.CreateCity(&city)
		if err != nil {
			t.Fatalf("CreateCity returned error %+v", err)
		}

		space := app.CitySpace{
			CityID:    city.ID,
			Order:     1,
			SpaceType: app.TraderID,
		}
		err = p.CreateCitySpace(&space)
		if err != nil {
			t.Fatalf("CreateCitySpace returned error %+v", err)
		}

		err = p.DeleteBoardByID(board.ID)
		if err != nil {
			t.Fatalf("DeleteBoardByID returned error %+v", err)
		}

		var cities []City
		err = tx.Find(&cities, "board_id = ?", board.ID).Error
		if err != nil {
			t.Fatalf("Finding cities by board id returned error %+v", err)
		}
		assert.ThatInt(len(cities)).IsEqualTo(0)

		var spaces []CitySpace
		err = tx.Find(&spaces, "city_id = ?", city.ID).Error
		assert.That(err).IsNil()
		assert.ThatInt(len(spaces)).IsEqualTo(0)

		var boards []Board
		err = tx.Find(&boards, board.ID).Error
		assert.That(err).IsNil()
		assert.ThatInt(len(boards)).IsEqualTo(0)
	})
}

func TestListCitiesByBoardId(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func (p app.BoardCrudRepository, tx *gorm.DB) {
		board := createTestBoard()
		createTestCityWithSpaces(board.ID)

		results, err := p.ListCitiesByBoardID(board.ID)
		if err != nil {
			t.Fatalf("ListCitiesByBoardID returned error: %+v", err)
		}
		assert.ThatInt(len(results)).IsEqualTo(1)
		assert.ThatInt(len(results[0].CitySpaces)).IsGreaterThan(0)
	})
}

func TestCreateCity(t *testing.T) {
	TempTransaction(func(p app.BoardCrudRepository, tx *gorm.DB) {
		board := createTestBoard()

		city := app.City{
			BoardID: board.ID,
			Name: "New City Name",
			Position: app.Position{
				X: 123,
				Y: 432,
			},
		}

		err := p.CreateCity(&city)
		if err != nil {
			t.Fatalf("CreateCity returned error: %+v", err)
		}

		var updatedCity City
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
	TempTransaction(func(p app.BoardCrudRepository, tx *gorm.DB) {
		board := createTestBoard()
		city := createTestCity(board.ID)

		err := p.UpdateCity(city.ID, func(city *app.City) (*app.City, error) {
			city.Name = "New City Name"
			city.Position.X = 123
			city.Position.Y = 432
			return city, nil
		})
		if err != nil {
			t.Fatalf("UpdateCity returned error: %+v", err)
		}

		var updatedCity City
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
	TempTransaction(func (p app.BoardCrudRepository, tx *gorm.DB) {
		board := createTestBoard()
		city := createTestCityWithSpaces(board.ID)

		spaces, err := p.GetCitySpacesByCityID(city.ID)
		assert.That(err).IsNil()

		assert.ThatInt(len(spaces)).IsEqualTo(len(city.CitySpaces))
	})
}

func TestSaveCitySpace(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func (r app.BoardCrudRepository, tx *gorm.DB) {
		board := createTestBoard()
		testCity := createTestCity(board.ID)

		space := app.CitySpace{
			CityID: testCity.ID,
			Order: 1,
			SpaceType: app.MerchantID,
			RequiredPrivilege: 2,
		}
		err := r.CreateCitySpace(&space)
		if err != nil {
			t.Fatalf("CreateCitySpace returned error: %+v", err)
		}

		city, err := r.GetCityByID(testCity.ID)
		if err != nil {
			t.Fatalf("GetCityById returned error: %+v", err)
		}

		space = city.CitySpaces[0]
		assert.That(space.SpaceType).IsEqualTo(app.MerchantID)
		assert.ThatInt(space.RequiredPrivilege).IsEqualTo(2)
		assert.ThatInt(space.Order).IsEqualTo(1)
	})
}

func TestDeleteCitySpaceByIDSucceeds(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func (p app.BoardCrudRepository, tx *gorm.DB) {
		board := createTestBoard()
		testCity := createTestCity(board.ID)

		space := app.CitySpace{
			CityID: testCity.ID,
			Order: 1,
			SpaceType: app.MerchantID,
			RequiredPrivilege: 2,
		}
		err := p.CreateCitySpace(&space)
		assert.That(err).IsNil()

		err = p.DeleteCitySpaceByID(space.ID)
		assert.That(err).IsNil()

		city, err := p.GetCityByID(testCity.ID)
		assert.That(err).IsNil()

		assert.ThatInt(len(city.CitySpaces)).IsEqualTo(0)
	})
}

var testBoardCounter = 0

func createTestBoard() *Board {
	board := Board{
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

func createTestCity(boardID uint) *City {
	city := City{
		BoardID: boardID,
		Name:    "Test City",
	}

	if err := db.Save(&city).Error; err != nil {
		panic(err)
	}

	return &city
}

func createTestCityWithSpaces(boardID uint) *City {
	city := createTestCity(boardID)

	city.CitySpaces = []CitySpace{
		{
			CityID:            city.ID,
			Order:             1,
			SpaceType:         app.TraderID,
			RequiredPrivilege: 1,
		},
		{
			CityID:            city.ID,
			Order:             2,
			SpaceType:         app.MerchantID,
			RequiredPrivilege: 2,
		},
		{
			CityID:            city.ID,
			Order:             3,
			SpaceType:         app.TraderID,
			RequiredPrivilege: 3,
		},
	}

	if err := db.Save(city.CitySpaces).Error; err != nil {
		panic(err)
	}

	return city
}
