package sqlc_board_crud_repository

import (
	"city-route-game/internal/app"
	"city-route-game/internal/sqlc"
	"context"
	"database/sql"
	"fmt"
)

func New(db *sql.DB) app.BoardCrudRepository {
	return &SqlcBoardCrudRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

type SqlcBoardCrudRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func (s SqlcBoardCrudRepository) transaction(ctx context.Context, fn func(q *sqlc.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = fn(s.q.WithTx(tx))

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			panic(rollbackErr)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return nil
}

func (s SqlcBoardCrudRepository) GetBoardByID(ctx context.Context, id app.ID) (*app.Board, error) {
	board, err := s.q.GetBoard(ctx, id)
	if err != nil {
		return nil, err
	}

	appBoard := newAppBoardFromGeneratedBoard(board)
	return &appBoard, nil
}

func (s SqlcBoardCrudRepository) CreateBoard(ctx context.Context, board *app.Board) error {
	params := sqlc.CreateBoardParams{
		Name:   board.Name,
		Width:  board.Width,
		Height: board.Height,
	}

	createdBoard, err := s.q.CreateBoard(ctx, params)
	if err != nil {
		return err
	}

	board.ID = createdBoard.ID
	board.Name = createdBoard.Name
	board.Width = createdBoard.Width
	board.Height = createdBoard.Height
	board.CreatedAt = createdBoard.CreatedAt
	board.UpdatedAt = createdBoard.UpdatedAt
	return nil
}

func (s SqlcBoardCrudRepository) UpdateBoard(
	ctx context.Context,
	id app.ID,
	updateFn func(board *app.Board) (*app.Board, error),
) (*app.Board, error) {

	var updatedAppBoard app.Board
	err := s.transaction(ctx, func(q *sqlc.Queries) error {
		originalBoard, err := q.GetBoard(ctx, id)
		if err != nil {
			return err
		}

		board := newAppBoardFromGeneratedBoard(originalBoard)
		userChanges, err := updateFn(&board)
		if err != nil {
			return err
		}

		params := sqlc.UpdateBoardParams{
			ID:     userChanges.ID,
			Name:   userChanges.Name,
			Width:  userChanges.Width,
			Height: userChanges.Height,
		}
		updatedBoard, err := q.UpdateBoard(ctx, params)
		if err != nil {
			return err
		}
		updatedAppBoard = newAppBoardFromGeneratedBoard(updatedBoard)
		return nil
	})
	return &updatedAppBoard, err
}

func (s SqlcBoardCrudRepository) ListBoards(ctx context.Context) ([]app.Board, error) {
	results, err := s.q.ListBoards(ctx)
	if err != nil {
		return nil, err
	}

	boards := make([]app.Board, 0, len(results))
	for _, row := range results {
		board := newAppBoardFromGeneratedBoard(row)
		boards = append(boards, board)
	}
	return boards, nil
}

func (s SqlcBoardCrudRepository) DeleteBoardByID(ctx context.Context, id app.ID) error {
	return s.transaction(ctx, func(q *sqlc.Queries) error {
		_, err := q.GetBoard(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return app.NewBoardNotFoundError(id)
			}
			return err
		}

		cityIDs, err := q.ListCityIDsByBoardID(ctx, id)
		if err != nil {
			return err
		}

		err = q.DeleteCitySpacesWhereCityIDIn(ctx, cityIDs)
		if err != nil {
			return err
		}

		err = q.DeleteMultipleCities(ctx, cityIDs)
		if err != nil {
			return err
		}

		err = q.DeleteBoard(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s SqlcBoardCrudRepository) ListCitiesByBoardID(ctx context.Context, boardID app.ID) ([]app.City, error) {
	var cities []app.City
	err := s.transaction(ctx, func(q *sqlc.Queries) error {
		sqlCities, err := q.ListCitiesByBoardID(ctx, boardID)
		if err != nil {
			return err
		}

		cityIDs := make([]int64, 0, len(sqlCities))
		for _, city := range sqlCities {
			cityIDs = append(cityIDs, city.ID)
		}

		citySpaces, err := q.ListCitySpacesByMultipleCities(ctx, cityIDs)
		if err != nil {
			return err
		}

		cities = make([]app.City, 0, len(sqlCities))
		for _, c := range sqlCities {
			city := newAppCityFromGeneratedCity(c)
			cities = append(cities, city)
		}

		cityIdx := 0
		for _, s := range citySpaces {
			space := newAppCitySpaceFromGeneratedCitySpace(s)
			for space.CityID != cities[cityIdx].ID && cityIdx < len(cities) {
				cityIdx++
			}
			if cityIdx == len(cities)-1 && cities[cityIdx].ID != space.CityID {
				panic(fmt.Errorf("reached end of cities, but did not find city to assign space %d to", space.ID))
			}

			cities[cityIdx].CitySpaces = append(cities[cityIdx].CitySpaces, space)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return cities, nil
}

func (s SqlcBoardCrudRepository) GetCityByID(ctx context.Context, id app.ID) (*app.City, error) {
	var generatedCity sqlc.City
	var generatedSpaces []sqlc.CitySpace

	err := s.transaction(ctx, func(q *sqlc.Queries) error {
		var err error
		generatedCity, err = s.q.GetCity(ctx, id)
		if err != nil {
			return err
		}

		generatedSpaces, err = s.q.ListCitySpacesByCityID(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	city := newAppCityFromGeneratedCity(generatedCity)
	city.CitySpaces = make([]app.CitySpace, 0, len(generatedSpaces))
	for _, generatedSpace := range generatedSpaces {
		space := newAppCitySpaceFromGeneratedCitySpace(generatedSpace)
		city.CitySpaces = append(city.CitySpaces, space)
	}

	return &city, nil
}

func (s SqlcBoardCrudRepository) CreateCity(ctx context.Context, city *app.City) error {
	params := sqlc.CreateCityParams{
		BoardID:        city.BoardID,
		Name:           city.Name,
		X:              city.Position.X,
		Y:              city.Position.Y,
		UpgradeOffered: city.UpgradeOffered,
		ImmediatePoint: city.ImmediatePoint,
	}

	var createdCity sqlc.City
	err := s.transaction(ctx, func(q *sqlc.Queries) error {
		_, err := q.GetBoard(ctx, params.BoardID)
		if err != nil {
			return err
		}

		createdCity, err = q.CreateCity(ctx, params)
		return err
	})
	if err != nil {
		return err
	}

	city.ID = createdCity.ID
	city.CreatedAt = createdCity.CreatedAt
	city.UpdatedAt = createdCity.UpdatedAt
	return nil
}

func (s SqlcBoardCrudRepository) UpdateCity(
	ctx context.Context,
	id app.ID,
	updateFn func(city *app.City) (*app.City, error),
) (*app.City, error) {

	var updatedAppCity app.City

	err := s.transaction(ctx, func(q *sqlc.Queries) error {
		originalCity, err := q.GetCity(ctx, id)
		if err != nil {
			return err
		}

		city := newAppCityFromGeneratedCity(originalCity)
		userChanges, err := updateFn(&city)
		if err != nil {
			return err
		}

		params := sqlc.UpdateCityParams{
			ID:             userChanges.ID,
			Name:           userChanges.Name,
			X:              userChanges.Position.X,
			Y:              userChanges.Position.Y,
			UpgradeOffered: userChanges.UpgradeOffered,
			ImmediatePoint: userChanges.ImmediatePoint,
		}
		updatedCity, err := q.UpdateCity(ctx, params)
		if err != nil {
			return err
		}
		updatedAppCity = newAppCityFromGeneratedCity(updatedCity)
		return nil
	})
	return &updatedAppCity, err
}

func (s SqlcBoardCrudRepository) DeleteCityByBoardIDAndCityID(ctx context.Context, boardID app.ID, cityID app.ID) error {
	return s.transaction(ctx, func(q *sqlc.Queries) error {
		_, err := q.GetBoard(ctx, boardID)
		if err != nil {
			return err
		}
		_, err = q.GetCity(ctx, cityID)
		if err != nil {
			return err
		}

		err = q.DeleteCitySpacesWhereCityIDIn(ctx, []int64{cityID})
		if err != nil {
			return err
		}

		return q.DeleteCity(ctx, cityID)
	})
}

func (s SqlcBoardCrudRepository) DeleteCityByID(ctx context.Context, id app.ID) error {
	return s.transaction(ctx, func(q *sqlc.Queries) error {
		_, err := q.GetCity(ctx, id)
		if err != nil {
			return err
		}

		err = q.DeleteCitySpacesWhereCityIDIn(ctx, []int64{id})
		if err != nil {
			return err
		}

		return q.DeleteCity(ctx, id)
	})
}

func (s SqlcBoardCrudRepository) CreateCitySpace(ctx context.Context, space *app.CitySpace) error {
	if space == nil {
		panic("space must not be nil")
	}

	params := sqlc.CreateCitySpaceParams{
		CityID:            space.CityID,
		Order:             space.Order,
		SpaceType:         space.SpaceType,
		RequiredPrivilege: space.RequiredPrivilege,
	}

	createdSpace, err := s.q.CreateCitySpace(ctx, params)
	if err != nil {
		return err
	}

	space.ID = createdSpace.ID
	space.CreatedAt = createdSpace.CreatedAt
	space.UpdatedAt = createdSpace.UpdatedAt
	return nil
}

func (s SqlcBoardCrudRepository) UpdateCitySpace(
	ctx context.Context,
	id app.ID,
	updateFn func(space *app.CitySpace) (*app.CitySpace, error),
	) (*app.CitySpace, error) {
	var updatedAppCitySpace app.CitySpace

	err := s.transaction(ctx, func(q *sqlc.Queries) error {
		originalCitySpace, err := q.GetCitySpaceByID(ctx, id)
		if err != nil {
			return err
		}

		space := newAppCitySpaceFromGeneratedCitySpace(originalCitySpace)
		userChanges, err := updateFn(&space)
		if err != nil {
			return err
		}

		params := sqlc.UpdateCitySpaceParams{
			Order:             userChanges.Order,
			SpaceType:         userChanges.SpaceType,
			RequiredPrivilege: userChanges.RequiredPrivilege,
		}
		updatedCitySpace, err := q.UpdateCitySpace(ctx, params)
		if err != nil {
			return err
		}
		updatedAppCitySpace = newAppCitySpaceFromGeneratedCitySpace(updatedCitySpace)
		return nil
	})
	return &updatedAppCitySpace, err
}

func (s SqlcBoardCrudRepository) GetCitySpacesByCityID(ctx context.Context, cityID app.ID) ([]app.CitySpace, error) {
	generatedCitySpaces, err := s.q.ListCitySpacesByCityID(ctx, cityID)
	if err != nil {
		return nil, err
	}

	citySpaces := make([]app.CitySpace, 0, len(generatedCitySpaces))
	for _, gs := range generatedCitySpaces {
		space := newAppCitySpaceFromGeneratedCitySpace(gs)
		citySpaces = append(citySpaces, space)
	}

	return citySpaces, nil
}

func (s SqlcBoardCrudRepository) DeleteCitySpaceByID(ctx context.Context, id app.ID) error {
	return s.transaction(ctx, func(q *sqlc.Queries) error {
		_, err := q.GetCitySpaceByID(ctx, id)
		if err != nil {
			return err
		}

		return q.DeleteCitySpace(ctx, id)
	})
}
