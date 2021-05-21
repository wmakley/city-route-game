package domain

// PersistenceProvider Service that is capable of loading, saving, and deleting domain objects
type PersistenceProvider interface {
	// Transaction start a new transaction, or reuse the current one
	Transaction(callback func(PersistenceProvider) error) error

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
}
