package sqlc_board_crud_repository

import (
	"city-route-game/internal/app"
	"city-route-game/internal/sqlc"
	"context"
	"errors"
	"fmt"
	"github.com/assertgo/assert"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
	"database/sql"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
	repo app.BoardCrudRepository
)

func TestMain(m *testing.M) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading dotenv: %+v", err)
	}

	databaseUrl := os.Getenv("TEST_DATABASE_URL")
	db, err = sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatalf("Error connecting to database: %+v", err)
	}

	repo = New(db)

	os.Exit(m.Run())
}

// Run tests within a transaction that is always rolled back
func TempTransaction(callback func(context.Context, app.BoardCrudRepository, *sql.Tx)) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}
	defer (func () {
		err := tx.Rollback()
		if err != nil {
			panic(err)
		}
	})()
	callback(ctx, repo, tx)
}

func TestListBoards(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func(ctx context.Context, p app.BoardCrudRepository, tx *sql.Tx) {
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
			err := p.CreateBoard(ctx, &board)
			if err != nil {
				t.Fatalf("CreateBoard failed: %+v", err)
			}
		}

		results, err := p.ListBoards(ctx)
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
	TempTransaction(func(ctx context.Context, p app.BoardCrudRepository, tx *sql.Tx) {
		board := &app.Board{
			Name:   "My Awesome Board",
			Width:  200,
			Height: 300,
		}

		err := p.CreateBoard(ctx, board)
		if err != nil {
			t.Fatalf("CreateBoard returned error: %+v", err)
		}

		boardID := board.ID
		board, err = p.GetBoardByID(ctx, boardID)
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
	queries := sqlc.New(db)
	beginCount, err := queries.CountBoards(context.Background())
	if err != nil {
		panic(err)
	}

	assert := assert.New(t)
	TempTransaction(func(ctx context.Context, r app.BoardCrudRepository, tx *sql.Tx) {
		board1 := app.Board{
			Name:   "My Awesome Board",
			Width:  200,
			Height: 300,
		}

		err := r.CreateBoard(ctx, &board1)
		if err != nil {
			t.Fatalf("CreateBoard returned error: %+v", err)
		}

		board2 := app.Board{
			Name:   "My Awesome Board",
			Width:  200,
			Height: 300,
		}
		err = r.CreateBoard(ctx, &board2)
		if err != app.ErrNameTaken {
			t.Errorf("Expected ErrNameTaken to be returned for duplicate board name, got: %+v", err)
		}
	})

	endCount, err := queries.CountBoards(context.Background())
	if err != nil {
		panic(err)
	}
	assert.That(endCount).IsEqualTo(beginCount)
}

func TestUpdateBoard(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func(ctx context.Context, r app.BoardCrudRepository, tx *sql.Tx) {
		board := &app.Board{
			Name:   "Original Name",
			Width:  200,
			Height: 300,
		}
		err := r.CreateBoard(ctx, board)
		if err != nil {
			t.Fatalf("CreateBoard returned error: %+v", err)
		}

		originalID := board.ID

		board, err = r.UpdateBoard(ctx, originalID, func(board *app.Board) (*app.Board, error) {
			board.Name = "New Name"
			board.Width = 123
			board.Height = 321
			return board, nil
		})
		if err != nil {
			t.Fatalf("UpdateBoard returned error: %+v", err)
		}

		board, err = r.GetBoardByID(ctx, board.ID)
		if err != nil {
			t.Fatalf("GetBoardByID returned error: %+v", err)
		}

		assert.That(board.ID).IsEqualTo(originalID)
		assert.ThatString(board.Name).IsEqualTo("New Name")
		assert.ThatInt(board.Width).IsEqualTo(123)
		assert.ThatInt(board.Height).IsEqualTo(321)

		dupe := app.Board{
			Name: "Imadupe",
		}
		err = r.CreateBoard(ctx, &dupe)
		if err != nil {
			t.Fatalf("CreateBoard returned error: %+v", err)
		}

		board, err = r.UpdateBoard(ctx, originalID, func(board *app.Board) (*app.Board, error) {
			board.Name = "Imadupe"
			return board, nil
		})
		if err != app.ErrNameTaken {
			t.Errorf("Expected updating name to be same as other board to fail (err: %+v)", err)
		}
	})
}

func TestDeleteBoardByIDDeletesNestedRecords(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func(ctx context.Context, r app.BoardCrudRepository, tx *sql.Tx) {
		queries := sqlc.New(tx)
		board := app.Board{
			Name: "Test Board",
		}
		err := r.CreateBoard(ctx, &board)
		if err != nil {
			t.Fatalf("CreateBoard returned error %+v", err)
		}

		city := app.City{
			BoardID: board.ID,
			Name:    "Test City",
		}
		err = r.CreateCity(ctx, &city)
		if err != nil {
			t.Fatalf("CreateCity returned error %+v", err)
		}

		space := app.CitySpace{
			CityID:    city.ID,
			Order:     1,
			SpaceType: app.TraderID,
		}
		err = r.CreateCitySpace(ctx, &space)
		if err != nil {
			t.Fatalf("CreateCitySpace returned error %+v", err)
		}

		err = r.DeleteBoardByID(ctx, board.ID)
		if err != nil {
			t.Fatalf("DeleteBoardByID returned error %+v", err)
		}

		cities, err := queries.ListCitiesByBoardID(ctx, board.ID)
		if err != nil {
			t.Fatalf("Finding cities by board id returned error %+v", err)
		}
		assert.ThatInt(len(cities)).IsEqualTo(0)

		spaces, err := queries.ListCitySpacesByCityID(ctx, city.ID)
		assert.That(err).IsNil()
		assert.ThatInt(len(spaces)).IsEqualTo(0)

		boards, err := queries.ListBoards(ctx)
		assert.That(err).IsNil()
		assert.ThatInt(len(boards)).IsEqualTo(0)
	})
}

func TestListCitiesByBoardId(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func(ctx context.Context, r app.BoardCrudRepository, tx *sql.Tx) {
		_, err := r.ListCitiesByBoardID(ctx, 1234)
		if !errors.Is(app.RecordNotFound{}, err) {
			t.Errorf("did not receive RecordNotFound error when board didn't exist, got: %+v", err)
		}

		board := createTestBoard(ctx, tx)
		createTestCityWithSpaces(ctx, tx, board.ID)

		results, err := r.ListCitiesByBoardID(ctx, board.ID)
		if err != nil {
			t.Fatalf("ListCitiesByBoardID returned error: %+v", err)
		}
		assert.ThatInt(len(results)).IsEqualTo(1)
		assert.ThatInt(len(results[0].CitySpaces)).IsGreaterThan(0)
	})
}

func TestCreateCity(t *testing.T) {
	TempTransaction(func(ctx context.Context, p app.BoardCrudRepository, tx *sql.Tx) {
		queries := sqlc.New(tx)
		board := createTestBoard(ctx, tx)

		city := app.City{
			BoardID: board.ID,
			Name:    "New City Name",
			Position: app.Position{
				X: 123,
				Y: 432,
			},
		}

		err := p.CreateCity(ctx, &city)
		if err != nil {
			t.Fatalf("CreateCity returned error: %+v", err)
		}

		updatedCity, err := queries.GetCity(ctx, city.ID)
		if err != nil {
			t.Fatalf("GetCity returned error: %+v", err)
		}

		if updatedCity.Name != "New City Name" {
			t.Error("City Name was not updated")
		}
		if updatedCity.X != 123 {
			t.Error("City Position X was not updated")
		}
		if updatedCity.Y != 432 {
			t.Error("City Position Y was not updated")
		}
	})
}

func TestUpdateCity(t *testing.T) {
	TempTransaction(func(ctx context.Context, p app.BoardCrudRepository, tx *sql.Tx) {
		board := createTestBoard(ctx, tx)
		city := createTestCity(ctx, tx, board.ID)

		id := city.ID
		_, err := p.UpdateCity(ctx, id, func(city *app.City) (*app.City, error) {
			city.Name = "New City Name"
			city.Position.X = 123
			city.Position.Y = 432
			return city, nil
		})
		if err != nil {
			t.Fatalf("UpdateCity returned error: %+v", err)
		}

		updatedCity, err := p.GetCityByID(ctx, id)
		if err != nil {
			t.Fatalf("error reloading city: %+v", err)
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

func TestDeleteCityByID(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func(ctx context.Context, p app.BoardCrudRepository, tx *sql.Tx) {
		board := createTestBoard(ctx, tx)
		city := createTestCityWithSpaces(ctx, tx, board.ID)
		spaces, err := p.GetCitySpacesByCityID(ctx, city.ID)
		assert.That(err).IsNil()
		assert.That(len(spaces)).IsEqualTo(3)

		err = p.DeleteCityByID(ctx, city.ID)
		if err != nil {
			t.Fatalf("DeleteCityByID returner error: %+v", err)
		}

		spaces, err = p.GetCitySpacesByCityID(ctx, city.ID)
		assert.That(err).IsNil()

		assert.ThatInt(len(spaces)).IsEqualTo(0)
	})
}

func TestGetCitySpacesByCityID(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func(ctx context.Context, p app.BoardCrudRepository, tx *sql.Tx) {
		board := createTestBoard(ctx, tx)
		city := createTestCityWithSpaces(ctx, tx, board.ID)

		spaces, err := p.GetCitySpacesByCityID(ctx, city.ID)
		assert.That(err).IsNil()

		assert.ThatInt(len(spaces)).IsEqualTo(3)
	})
}

func TestSaveCitySpace(t *testing.T) {
	assert := assert.New(t)
	TempTransaction(func(ctx context.Context, r app.BoardCrudRepository, tx *sql.Tx) {
		board := createTestBoard(ctx, tx)
		testCity := createTestCity(ctx, tx, board.ID)

		space := app.CitySpace{
			CityID:            testCity.ID,
			Order:             1,
			SpaceType:         app.MerchantID,
			RequiredPrivilege: 2,
		}
		err := r.CreateCitySpace(ctx, &space)
		if err != nil {
			t.Fatalf("CreateCitySpace returned error: %+v", err)
		}

		city, err := r.GetCityByID(ctx, testCity.ID)
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
	TempTransaction(func(ctx context.Context, p app.BoardCrudRepository, tx *sql.Tx) {
		board := createTestBoard(ctx, tx)
		testCity := createTestCity(ctx, tx, board.ID)

		space := app.CitySpace{
			CityID:            testCity.ID,
			Order:             1,
			SpaceType:         app.MerchantID,
			RequiredPrivilege: 2,
		}
		err := p.CreateCitySpace(ctx, &space)
		assert.That(err).IsNil()

		err = p.DeleteCitySpaceByID(ctx, space.ID)
		assert.That(err).IsNil()

		city, err := p.GetCityByID(ctx, testCity.ID)
		assert.That(err).IsNil()

		assert.ThatInt(len(city.CitySpaces)).IsEqualTo(0)
	})
}

var testBoardCounter = 0

func createTestBoard(ctx context.Context, tx *sql.Tx) sqlc.Board {
	queries := sqlc.New(tx)
	params := sqlc.CreateBoardParams{
		Name:   fmt.Sprintf("Test Board %d", testBoardCounter),
		Width:  10,
		Height: 20,
	}
	testBoardCounter++
	board, err := queries.CreateBoard(ctx, params)
	if err != nil {
		panic(err)
	}
	return board
}

func createTestCity(ctx context.Context, tx *sql.Tx, boardID int64) sqlc.City {
	queries := sqlc.New(tx)
	params := sqlc.CreateCityParams{
		BoardID: boardID,
		Name:    "Test City",
	}

	city, err := queries.CreateCity(ctx, params)
	if err != nil {
		panic(err)
	}

	return city
}

func createTestCityWithSpaces(ctx context.Context, tx *sql.Tx, boardID int64) sqlc.City {
	city := createTestCity(ctx, tx, boardID)
	queries := sqlc.New(tx)

	params := []sqlc.CreateCitySpaceParams{
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

	for _, p := range params {
		_, err := queries.CreateCitySpace(ctx, p)
		if err != nil {
			panic(err)
		}
	}

	return city
}
