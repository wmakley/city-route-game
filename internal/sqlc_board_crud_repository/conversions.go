package sqlc_board_crud_repository

import (
	"city-route-game/internal/app"
	"city-route-game/internal/sqlc"
)

func newAppBoardFromGeneratedBoard(board sqlc.Board) app.Board {
	return app.Board{
		Model: app.Model{
			ID:        board.ID,
			CreatedAt: board.CreatedAt,
			UpdatedAt: board.UpdatedAt,
		},
		Name:   board.Name,
		Width:  board.Width,
		Height: board.Height,
		Cities: nil,
	}
}

func newAppCityFromGeneratedCity(city sqlc.City) app.City {
	return app.City{
		Model: app.Model{
			ID:        city.ID,
			CreatedAt: city.CreatedAt,
			UpdatedAt: city.UpdatedAt,
		},
		BoardID: city.BoardID,
		Name:    city.Name,
		Position: app.Position{
			X: city.X,
			Y: city.Y,
		},
		CitySpaces: nil,
	}
}

func newAppCitySpaceFromGeneratedCitySpace(citySpace sqlc.CitySpace) app.CitySpace {
	return app.CitySpace{
		Model:             app.Model{
			ID: citySpace.ID,
			CreatedAt: citySpace.CreatedAt,
			UpdatedAt: citySpace.UpdatedAt,
		},
		CityID:            citySpace.CityID,
		Order:             citySpace.Order,
		SpaceType:         citySpace.SpaceType,
		RequiredPrivilege: citySpace.RequiredPrivilege,
	}
}
