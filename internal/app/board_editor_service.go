package app

import (
	"errors"
)

type BoardEditorService interface {
	FindAll() ([]Board, error)
	FindByID(id ID) (*Board, error)
	CreateBoard(form *BoardNameForm) (*Board, error)
	UpdateName(id ID, form *BoardNameForm) (*Board, error)
	UpdateDimensions(id ID, form *BoardDimensionsForm) (*Board, error)
	DeleteByID(id ID) error
}

func NewBoardEditorService(boardCrudRepository BoardCrudRepository) BoardEditorService {
	return &boardEditorService{
		repo: boardCrudRepository,
	}
}

type boardEditorService struct {
	repo BoardCrudRepository
}

func (s boardEditorService)FindAll() ([]Board, error) {
	return s.repo.ListBoards()
}

func (s boardEditorService)FindByID(id ID) (*Board, error) {
	board, err := s.repo.GetBoardByID(id)
	if err != nil {
		return nil, err
	}
	return board, nil
}

func (s boardEditorService)CreateBoard(form *BoardNameForm) (*Board, error) {
	form.NormalizeInputs()

	if !form.IsValid() {
		return nil, ErrInvalidForm
	}

	board := Board{
		Name:   form.Name,
		Width:  800,
		Height: 600,
	}

	err := s.repo.CreateBoard(&board)
	if err != nil {
	if errors.Is(ErrNameTaken, err) {
			form.AddError("Name", "is already taken")
			return nil, ErrInvalidForm
		}

		return nil, err
	}

	return &board, nil
}

func (s boardEditorService)UpdateDimensions(id ID, form *BoardDimensionsForm) (*Board, error) {
	return nil, nil
}

func (s boardEditorService)UpdateName(id ID, form *BoardNameForm) (*Board, error) {
	form.NormalizeInputs()

	if !form.IsValid() {
		return nil, ErrInvalidForm
	}

	return s.repo.UpdateBoard(form.ID, func (board *Board) (*Board, error) {
		board.Name = form.Name
		return board, nil
	})
}

func (s boardEditorService)DeleteByID(id ID) error {
	return nil
}
