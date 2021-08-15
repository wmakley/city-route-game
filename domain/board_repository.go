package domain

// BoardRepository Service that is capable of loading, saving, and deleting boards and board parts
type BoardRepository interface {
	GetBoardByID(id interface{}) (*Board, error)
	CreateBoard(board *Board) error
	UpdateBoard(id interface{}, updateFn func (board *Board) (*Board, error)) error
	ListBoards() ([]Board, error)
	DeleteBoardByID(id interface{}) error
	//BoardExistsWithName(name string) (bool, error)
	//BoardExistsWithNameAndIdNot(name string, idNot interface{}) (bool, error)

	ListCitiesByBoardID(boardID interface{}) ([]City, error)
	GetCityByID(id interface{}) (*City, error)
	CreateCity(city *City) error
	UpdateCity(id interface{}, updateFn func (city *City) (*City, error)) error
	DeleteCityByBoardIDAndCityID(boardID interface{}, cityID interface{}) error
	DeleteCity(city *City) error

	CreateCitySpace(*CitySpace) error
	UpdateCitySpace(id interface{}, updateFn func (space *CitySpace) (*CitySpace, error)) error
	GetCitySpacesByCityID(cityID interface{}) ([]CitySpace, error)
	DeleteCitySpaceByID(id interface{}) error
}
