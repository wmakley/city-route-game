package app

import (
	"errors"
	"testing"
)

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
	if err == nil {
		t.Error("CreateBoard with duplicate name should have returned an error")
	} else if !errors.Is(ErrInvalidForm, err) {
		t.Errorf("CreateBoard should have returned ErrInvalidForm, was: %+v", err)
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
