package app

import (
	"context"
	"errors"
	"strings"
)

type BoardEditorService interface {
	FindAll(ctx context.Context) ([]Board, error)
	FindByID(ctx context.Context, id string) (*Board, error)
	CreateBoard(ctx context.Context, form *CreateBoardForm) (*Board, error)
	Update(ctx context.Context, id string, form *UpdateBoardForm) (*Board, error)
	UpdateName(ctx context.Context, id string, form *UpdateBoardForm) (*Board, error)
	UpdateDimensions(ctx context.Context, id string, form *UpdateBoardForm) (*Board, error)
	DeleteByID(ctx context.Context, id string) error

	ListCitiesByBoardID(ctx context.Context, boardID string) ([]City, error)
	CreateCity(ctx context.Context, boardID string, form *CityForm) (*City, error)
	UpdateCity(ctx context.Context, id string, form *CityForm) (*City, error)
	DeleteCity(ctx context.Context, id string) error
}

func NewBoardEditorService(boardCrudRepository BoardCrudRepository) BoardEditorService {
	return &boardEditorService{
		repo: boardCrudRepository,
	}
}

type boardEditorService struct {
	repo BoardCrudRepository
}

func (s boardEditorService)FindAll(ctx context.Context) ([]Board, error) {
	return s.repo.ListBoards(ctx)
}

func (s boardEditorService)FindByID(ctx context.Context, rawId string) (*Board, error) {
	id, err := NewIDFromString(rawId)
	if err != nil {
		return nil, err
	}

	board, err := s.repo.GetBoardByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return board, nil
}

func (s boardEditorService)CreateBoard(ctx context.Context, form *CreateBoardForm) (*Board, error) {
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

	err := s.repo.CreateBoard(ctx, &board)
	if err != nil {
	if errors.Is(ErrNameTaken, err) {
			form.AddError("Name", "is already taken")
			return nil, ErrInvalidForm
		}

		return nil, err
	}

	return &board, nil
}

func (s boardEditorService)UpdateDimensions(ctx context.Context, rawId string, form *UpdateBoardForm) (*Board, error) {
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

	updatedBoard, err := s.repo.UpdateBoard(ctx, id, func (board *Board) (*Board, error) {
		board.Width = form.Width
		board.Height = form.Height
		return board, nil
	})
	if err != nil {
		return nil, err
	}

	return updatedBoard, nil
}

func (s boardEditorService)UpdateName(ctx context.Context, rawId string, form *UpdateBoardForm) (*Board, error) {
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

	updatedBoard, err := s.repo.UpdateBoard(ctx, id, func (board *Board) (*Board, error) {
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

func (s boardEditorService)Update(ctx context.Context, rawId string, form *UpdateBoardForm) (*Board, error) {
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

	updatedBoard, err := s.repo.UpdateBoard(ctx, id, func (board *Board) (*Board, error) {
		board.Name = form.Name
		board.Width = form.Width
		board.Height = form.Height
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

func (s boardEditorService)DeleteByID(ctx context.Context, rawId string) error {
	id, err := NewIDFromString(rawId)
	if err != nil {
		return err
	}

	return s.repo.DeleteBoardByID(ctx, id)
}

func (s boardEditorService)ListCitiesByBoardID(ctx context.Context, boardID string) ([]City, error) {
	id, err := NewIDFromString(boardID)
	if err != nil {
		return nil, err
	}

	return s.repo.ListCitiesByBoardID(ctx, id)
}

func (s boardEditorService)CreateCity(ctx context.Context, boardID string, form *CityForm) (*City, error) {
	parsedBoardID, err := NewIDFromString(boardID)
	if err != nil {
		return nil, err
	}

	form.NormalizeInputs()

	if !form.IsValid() {
		return nil, ErrInvalidForm
	}

	city := City{
		Model:      Model{},
		BoardID:    parsedBoardID,
		Name:       form.Name,
		Position:   form.Position,
		CitySpaces: nil,
	}

	err = s.repo.CreateCity(ctx, &city)
	return &city, err
}

func (s boardEditorService)UpdateCity(ctx context.Context, id string, form *CityForm) (*City, error) {
	parsedID, err := NewIDFromString(id)
	if err != nil {
		return nil, err
	}

	form.NormalizeInputs()

	if !form.IsValid() {
		return nil, ErrInvalidForm
	}

	updatedCity, err := s.repo.UpdateCity(ctx, parsedID, func(city *City) (*City, error) {
		city.Name = form.Name
		city.Position.X	= form.Position.X
		city.Position.Y = form.Position.Y
		return city, nil
	})
	if err != nil {
		return nil, err
	}
	return updatedCity, nil
}

func (s boardEditorService)DeleteCity(ctx context.Context, id string) error {
	parsedID, err := NewIDFromString(id)
	if err != nil {
		return err
	}

	return s.repo.DeleteCityByID(ctx, parsedID)
}
