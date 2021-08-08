package domain

// BoardRepository Service that is capable of loading, saving, and deleting boards and board pieces
type BoardRepository interface {
	GetBoardByID(id interface{}) (*Board, error)
	SaveBoard(board *Board) error
	ListBoards() ([]Board, error)
	DeleteBoardByID(id interface{}) error
	BoardExistsWithName(name string) (bool, error)
	BoardExistsWithNameAndIdNot(name string, idNot interface{}) (bool, error)

	ListCitiesByBoardID(boardID interface{}) ([]City, error)
	GetCityByID(id interface{}) (*City, error)
	SaveCity(city *City) error
	DeleteCityByBoardIDAndCityID(boardID interface{}, cityID interface{}) error
	DeleteCity(city *City) error

	SaveCitySpace(*CitySpace) error
	GetCitySpacesByCityID(cityID interface{}) ([]CitySpace, error)
	DeleteCitySpaceByID(id interface{}) error
}
