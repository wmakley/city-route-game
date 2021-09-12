package app

// BoardCrudRepository Repository that is capable of loading, saving, and deleting boards and board parts
type BoardCrudRepository interface {
	GetBoardByID(id ID) (*Board, error)
	CreateBoard(board *Board) error
	UpdateBoard(id ID, updateFn func (board *Board) (*Board, error)) (*Board, error)
	ListBoards() ([]Board, error)
	DeleteBoardByID(id ID) error
	//BoardExistsWithName(name string) (bool, error)
	//BoardExistsWithNameAndIdNot(name string, idNot interface{}) (bool, error)

	ListCitiesByBoardID(boardID ID) ([]City, error)
	GetCityByID(id ID) (*City, error)
	CreateCity(city *City) error
	UpdateCity(id ID, updateFn func (city *City) (*City, error)) error
	DeleteCityByBoardIDAndCityID(boardID ID, cityID ID) error
	DeleteCity(city *City) error

	CreateCitySpace(*CitySpace) error
	UpdateCitySpace(id ID, updateFn func (space *CitySpace) (*CitySpace, error)) error
	GetCitySpacesByCityID(cityID ID) ([]CitySpace, error)
	DeleteCitySpaceByID(id ID) error
}
