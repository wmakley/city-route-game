package app

import "context"

// BoardCrudRepository Repository that is capable of loading, saving, and deleting boards and board parts
type BoardCrudRepository interface {
	GetBoardByID(ctx context.Context, id ID) (*Board, error)
	CreateBoard(ctx context.Context, board *Board) error
	UpdateBoard(ctx context.Context, id ID, updateFn func (board *Board) (*Board, error)) (*Board, error)
	ListBoards(ctx context.Context) ([]Board, error)
	DeleteBoardByID(ctx context.Context, id ID) error
	//BoardExistsWithName(name string) (bool, error)
	//BoardExistsWithNameAndIdNot(name string, idNot interface{}) (bool, error)

	ListCitiesByBoardID(ctx context.Context, boardID ID) ([]City, error)
	GetCityByID(ctx context.Context, id ID) (*City, error)
	CreateCity(ctx context.Context, city *City) error
	UpdateCity(ctx context.Context, id ID, updateFn func (city *City) (*City, error)) (*City, error)
	DeleteCityByBoardIDAndCityID(ctx context.Context, boardID ID, cityID ID) error
	DeleteCityByID(ctx context.Context, id ID) error

	CreateCitySpace(context.Context, *CitySpace) error
	UpdateCitySpace(ctx context.Context, id ID, updateFn func (space *CitySpace) (*CitySpace, error)) (*CitySpace, error)
	GetCitySpacesByCityID(ctx context.Context, cityID ID) ([]CitySpace, error)
	DeleteCitySpaceByID(ctx context.Context, id ID) error
}
