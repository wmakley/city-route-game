package app

import (
	"errors"
	"github.com/assertgo/assert"
	"testing"
	"time"
)

func TestFindAll(t *testing.T) {
	repo := fakeBoardCrudRepository{}
	service := NewBoardEditorService(&repo)

	results, err := service.FindAll()
	if err != nil {
		t.Fatalf("FindAll returned error: %+v", err)
	}
	if len(results) != 0 {
		t.Error("Results slice should have been empty with no boards created yet")
	}

	now := time.Now()
	repo.Boards = []Board{
		{
			Model:  Model{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Name:   "Board 1",
			Width:  10,
			Height: 20,
			Cities: nil,
		},
		{
			Model:  Model{
				ID:        2,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Name:   "Board 2",
			Width:  30,
			Height: 40,
			Cities: nil,
		},
	}

	results, err = service.FindAll()
	if err != nil {
		t.Fatalf("FindAll returned error: %+v", err)
	}
	if len(results) != 2 {
		t.Errorf("Results length should have been 2 (was %d)", len(results))
	}

	assert := assert.New(t)
	assert.That(results[0].ID).IsEqualTo(ID(1))
	assert.That(results[0].CreatedAt).IsEqualTo(now)
	assert.That(results[0].UpdatedAt).IsEqualTo(now)
	assert.That(results[0].Name).IsEqualTo("Board 1")
	assert.That(results[0].Width).IsEqualTo(10)
	assert.That(results[0].Height).IsEqualTo(20)
	assert.ThatInt(len(results[0].Cities)).IsEqualTo(0)

	assert.That(results[1].ID).IsEqualTo(ID(2))
	assert.That(results[1].CreatedAt).IsEqualTo(now)
	assert.That(results[1].UpdatedAt).IsEqualTo(now)
	assert.That(results[1].Name).IsEqualTo("Board 2")
	assert.That(results[1].Width).IsEqualTo(30)
	assert.That(results[1].Height).IsEqualTo(40)
	assert.ThatInt(len(results[1].Cities)).IsEqualTo(0)
}

func FindByID(t *testing.T) {
	repo := fakeBoardCrudRepository{}
	service := NewBoardEditorService(&repo)

	_, err := service.FindByID("1")
	if !errors.Is(RecordNotFound{}, err) {
		t.Errorf("FindByID with no board should have returned error RecordNotFound, but returned: %+v", err)
	}

	now := time.Now()
	repo.Boards = []Board{
		{
			Model:  Model{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Name:   "Board 1",
			Width:  10,
			Height: 20,
			Cities: []City{
				{
					Model:      Model{
						ID:        1,
						CreatedAt: now,
						UpdatedAt: now,
					},
					BoardID:    1,
					Name:       "City 1",
					Position:   Position{
						X: 111,
						Y: 222,
					},
					CitySpaces: []CitySpace{
						{
							Model:             Model{
								ID:        1,
								CreatedAt: now,
								UpdatedAt: now,
							},
							CityID:            1,
							Order:             1,
							SpaceType:         TraderID,
							RequiredPrivilege: 1,
						},
					},
				},
			},
		},
	}

	var result *Board
	result, err = service.FindByID("1")
	if err != nil {
		t.Fatalf("FindByID returned error: %+v", err)
	}

	assert := assert.New(t)
	assert.That(result.ID).IsEqualTo(ID(1))
	assert.That(result.CreatedAt).IsEqualTo(now)
	assert.That(result.UpdatedAt).IsEqualTo(now)
	assert.That(result.Name).IsEqualTo("Board 1")
	assert.That(result.Width).IsEqualTo(10)
	assert.That(result.Height).IsEqualTo(20)
	assert.ThatInt(len(result.Cities)).IsEqualTo(1)

	city := result.Cities[0]
	assert.That(city.ID).IsEqualTo(ID(1))
	assert.That(city.CreatedAt).IsEqualTo(now)
	assert.That(city.UpdatedAt).IsEqualTo(now)
	assert.That(city.BoardID).IsEqualTo(ID(1))
	assert.That(city.Name).IsEqualTo("City 1")
	assert.That(city.Position.X).IsEqualTo(111)
	assert.That(city.Position.Y).IsEqualTo(222)
	assert.ThatInt(len(city.CitySpaces)).IsEqualTo(1)

	space := city.CitySpaces[0]
	assert.That(space.ID).IsEqualTo(ID(1))
	assert.That(space.CreatedAt).IsEqualTo(now)
	assert.That(space.UpdatedAt).IsEqualTo(now)
	assert.That(space.Order).IsEqualTo(1)
	assert.That(space.SpaceType).IsEqualTo(TraderID)
	assert.That(space.CityID).IsEqualTo(ID(1))
	assert.That(space.RequiredPrivilege).IsEqualTo(1)
}

func TestCreateBoard(t *testing.T) {
	repo := fakeBoardCrudRepository{}
	service := NewBoardEditorService(&repo)

	form := NewCreateBoardForm()
	form.Name = "Test Name"

	board, err := service.CreateBoard(&form)
	if err != nil {
		t.Errorf("CreateBoard with valid name returned error: %+v", err)
	}
	if board.Width <= 0 {
		t.Errorf("Expected non-zero default Width (was %d)", board.Width)
	}
	if board.Height <= 0 {
		t.Errorf("Exected non-zeri default Height (was %d)", board.Height)
	}

	form.Name = ""
	_, err = service.CreateBoard(&form)
	if err == nil {
		t.Error("CreateBoard with blank name should have returned an error")
	} else if !errors.Is(ErrInvalidForm, err) {
		t.Errorf("CreateBoard should have returned ErrInvalidForm, was: %+v", err)
	}

	form.Name = "   "
	_, err = service.CreateBoard(&form)
	if err == nil {
		t.Error("CreateBoard with blank name should have returned an error")
	} else if !errors.Is(ErrInvalidForm, err) {
		t.Errorf("CreateBoard should have returned ErrInvalidForm, was: %+v", err)
	}

	repo.ErrorResult = ErrNameTaken
	form.Name = "Duplicate Name"
	_, err = service.CreateBoard(&form)
	if !errors.Is(ErrInvalidForm, err) {
		t.Errorf("CreateBoard should have returned ErrInvalidForm, was: %+v", err)
	}
	_, ok := form.Errors["Name"]
	if !ok {
		t.Error("No error for 'Name' was found in form")
	}
}

func TestUpdateName(t *testing.T) {
	now := time.Now()
	repo := fakeBoardCrudRepository{
		Boards: []Board{
			{
				Model:  Model{
					ID: 1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Name:   "Original Name",
				Width:  10,
				Height: 20,
				Cities: nil,
			},
		},
	}
	service := NewBoardEditorService(&repo)

	form := NewBoardNameForm(&repo.Boards[0])
	form.Name = "Test Name"

	updatedBoard, err := service.UpdateName("1", &form)
	if err != nil {
		t.Errorf("UpdateName with valid name returned error: %+v", err)
	}
	if updatedBoard.Name != "Test Name" {
		t.Error("new name was not set")
	}

	form.Name = ""
	_, err = service.UpdateName("1", &form)
	if err == nil {
		t.Error("UpdateName with blank name should have returned an error")
	} else if !errors.Is(ErrInvalidForm, err) {
		t.Errorf("UpdateName should have returned ErrInvalidForm, was: %+v", err)
	}

	form.Name = "   "
	_, err = service.UpdateName("1", &form)
	if err == nil {
		t.Error("UpdateName with blank name should have returned an error")
	} else if !errors.Is(ErrInvalidForm, err) {
		t.Errorf("UpdateName should have returned ErrInvalidForm, was: %+v", err)
	}

	repo.ErrorResult = ErrNameTaken
	form.Name = "Duplicate Name"
	_, err = service.UpdateName("1", &form)
	if !errors.Is(ErrInvalidForm, err) {
		t.Errorf("UpdateBoard should have returned ErrInvalidForm, was: %+v", err)
	}
	_, ok := form.Errors["Name"]
	if !ok {
		t.Error("No error for 'Name' was found in form")
	}
}

func TestUpdateDimensions(t *testing.T) {
	assert := assert.New(t)
	now := time.Now()
	repo := fakeBoardCrudRepository{
		Boards: []Board{
			{
				Model:  Model{
					ID: 1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Name:   "Board 1",
				Width:  10,
				Height: 20,
				Cities: nil,
			},
		},
	}
	service := NewBoardEditorService(&repo)

	form := NewBoardDimensionsForm(&repo.Boards[0])
	form.Width = 123
	form.Height = 456

	updatedBoard, err := service.UpdateDimensions("1", &form)
	if err != nil {
		t.Fatalf("UpdateDimensions with valid input returned error: %+v", err)
	}
	assert.That(updatedBoard.Name).IsEqualTo("Board 1")
	assert.That(updatedBoard.Width).IsEqualTo(123)
	assert.That(updatedBoard.Height).IsEqualTo(456)

	form.Width = -1
	form.Height = 0
	_, err = service.UpdateDimensions("1", &form)
	if err == nil {
		t.Fatalf("UpdateDimensions with negative width should have returned an error")
	} else if !errors.Is(ErrInvalidForm, err) {
		t.Errorf("UpdateDimensions with negative width should have returned ErrInvalidForm, was: %+v", err)
	}

	form.Width = 0
	form.Height = -1
	_, err = service.UpdateDimensions("1", &form)
	if err == nil {
		t.Fatalf("UpdateDimensions with negative height should have returned an error")
	} else if !errors.Is(ErrInvalidForm, err) {
		t.Errorf("UpdateDimensions with negative height should have returned ErrInvalidForm, was: %+v", err)
	}
}

func TestDeleteByID(t *testing.T) {
	repo := fakeBoardCrudRepository{}
	service := NewBoardEditorService(&repo)

	err := service.DeleteByID("1")
	if err != nil {
		t.Error("DeleteByID should have returned nil error when repo returned success")
	}

	repo.ErrorResult = NewBoardNotFoundError(1)
	err = service.DeleteByID("1")
	if !errors.Is(RecordNotFound{}, err) {
		t.Errorf("DeleteByID when board doesn't exist should have returned RecordNotFound, was: %+v", err)
	}
}

type fakeBoardCrudRepository struct {
	Boards []Board
	SingletonCityResult *City
	MultipleCityResult []City
	CitySpaces []CitySpace
	ErrorResult error
}

func (r fakeBoardCrudRepository)GetBoardByID(id ID) (*Board, error) {
	return &r.Boards[0], r.ErrorResult
}
func (r fakeBoardCrudRepository)CreateBoard(board *Board) error {
	return r.ErrorResult
}
func (r fakeBoardCrudRepository)UpdateBoard(id ID, updateFn func (board *Board) (*Board, error)) (*Board, error) {
	result, err := updateFn(&r.Boards[0])
	if err != nil {
		return nil, err
	}
	return result, r.ErrorResult
}
func (r fakeBoardCrudRepository)ListBoards() ([]Board, error) {
	return r.Boards, r.ErrorResult
}
func (r fakeBoardCrudRepository)DeleteBoardByID(id ID) error {
	return r.ErrorResult
}

func (r fakeBoardCrudRepository)ListCitiesByBoardID(boardID ID) ([]City, error) {
	return r.MultipleCityResult, r.ErrorResult
}
func (r fakeBoardCrudRepository)GetCityByID(id ID) (*City, error) {
	return r.SingletonCityResult, r.ErrorResult
}
func (r fakeBoardCrudRepository)CreateCity(city *City) error {
	return r.ErrorResult
}
func (r fakeBoardCrudRepository)UpdateCity(id ID, updateFn func (city *City) (*City, error)) error {
	_, err := updateFn(r.SingletonCityResult)
	if err != nil {
		return err
	}
	return r.ErrorResult
}
func (r fakeBoardCrudRepository)DeleteCityByBoardIDAndCityID(boardID ID, cityID ID) error {
	return r.ErrorResult
}
func (r fakeBoardCrudRepository)DeleteCity(city *City) error {
	return r.ErrorResult
}

func (r fakeBoardCrudRepository)CreateCitySpace(*CitySpace) error {
	return r.ErrorResult
}
func (r fakeBoardCrudRepository)UpdateCitySpace(id ID, updateFn func (space *CitySpace) (*CitySpace, error)) error {
	_, err := updateFn(&r.CitySpaces[0])
	if err != nil {
		return err
	}
	return r.ErrorResult
}
func (r fakeBoardCrudRepository)GetCitySpacesByCityID(cityID ID) ([]CitySpace, error) {
	return r.CitySpaces, r.ErrorResult
}
func (r fakeBoardCrudRepository)DeleteCitySpaceByID(id ID) error{
	return r.ErrorResult
}
