package app

import (
	"errors"
	"strings"
)

type BoardEditorService interface {
	FindAll() ([]Board, error)
	FindByID(id string) (*Board, error)
	CreateBoard(form *CreateBoardForm) (*Board, error)
	Update(id string, form *UpdateBoardForm) (*Board, error)
	UpdateName(id string, form *UpdateBoardForm) (*Board, error)
	UpdateDimensions(id string, form *UpdateBoardForm) (*Board, error)
	DeleteByID(id string) error
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

func (s boardEditorService)FindByID(rawId string) (*Board, error) {
	id, err := NewIDFromString(rawId)
	if err != nil {
		return nil, err
	}

	board, err := s.repo.GetBoardByID(id)
	if err != nil {
		return nil, err
	}
	return board, nil
}

func (s boardEditorService)CreateBoard(form *CreateBoardForm) (*Board, error) {
	form.Name = strings.TrimSpace(form.Name)

	if len(form.Name) == 0 {
		form.AddError("Name", "must not be blank")
	} else if len(form.Name) > 100 {
		form.AddError("Name", "is too long; must be 100 characters or less")
	}

	if form.HasError() {
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

func (s boardEditorService)UpdateDimensions(rawId string, form *UpdateBoardForm) (*Board, error) {
	id, err := NewIDFromString(rawId)
	if err != nil {
		return nil, err
	}

	if form.Width < 0 {
		form.AddError("Width", "must be greater than or equal to zero")
	}
	if form.Height < 0 {
		form.AddError("Height", "must be greater than or equal to zero")
	}

	if form.HasError() {
		return nil, ErrInvalidForm
	}

	updatedBoard, err := s.repo.UpdateBoard(id, func (board *Board) (*Board, error) {
		board.Width = form.Width
		board.Height = form.Height
		return board, nil
	})
	if err != nil {
		return nil, err
	}

	return updatedBoard, nil
}

func (s boardEditorService)UpdateName(rawId string, form *UpdateBoardForm) (*Board, error) {
	id, err := NewIDFromString(rawId)
	if err != nil {
		return nil, err
	}

	form.Name = strings.TrimSpace(form.Name)

	if len(form.Name) == 0 {
		form.AddError("Name", "must not be blank")
	} else if len(form.Name) > 100 {
		form.AddError("Name", "is too long; must be 100 characters or less")
	}

	if form.HasError() {
		return nil, ErrInvalidForm
	}

	updatedBoard, err := s.repo.UpdateBoard(id, func (board *Board) (*Board, error) {
		board.Name = form.Name
		return board, nil
	})
	if err != nil {
		if errors.Is(ErrNameTaken, err) {
			form.AddError("Name", "name is already taken")
			return nil, ErrInvalidForm
		}
		return nil, err
	}

	return updatedBoard, nil
}

func (s boardEditorService)Update(rawId string, form *UpdateBoardForm) (*Board, error) {
	id, err := NewIDFromString(rawId)
	if err != nil {
		return nil, err
	}

	form.Name = strings.TrimSpace(form.Name)

	if len(form.Name) == 0 {
		form.AddError("Name", "must not be blank")
	} else if len(form.Name) > 100 {
		form.AddError("Name", "is too long; must be 100 characters or less")
	}

	if form.Width < 0 {
		form.AddError("Width", "must be greater than or equal to zero")
	}
	if form.Height < 0 {
		form.AddError("Height", "must be greater than or equal to zero")
	}

	if form.HasError() {
		return nil, ErrInvalidForm
	}

	updatedBoard, err := s.repo.UpdateBoard(id, func (board *Board) (*Board, error) {
		board.Name = form.Name
		board.Width = form.Width
		board.Name = form.Name
		return board, nil
	})
	if err != nil {
		if errors.Is(ErrNameTaken, err) {
			form.AddError("Name", "name is already taken")
			return nil, ErrInvalidForm
		}
		return nil, err
	}

	return updatedBoard, nil
}

func (s boardEditorService)DeleteByID(rawId string) error {
	id, err := NewIDFromString(rawId)
	if err != nil {
		return err
	}

	return s.repo.DeleteBoardByID(id)
}
