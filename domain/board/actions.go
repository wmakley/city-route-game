package board

import (
	"city-route-game/domain"
)

var repo Repository

func Init(repo_ Repository) {
	repo = repo_
}

func FindAll() ([]domain.Board, error) {
	results, err := repo.ListBoards()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func FindByID(id interface{}) (board *domain.Board, err error) {
	board, err = repo.GetBoardByID(id)
	return
}

func Create(form *Form) (*domain.Board, error) {
	form.NormalizeInputs()

	var board domain.Board
	err := repo.Transaction(func(tx Repository) error {
		if !form.IsValid(tx) {
			return domain.ErrInvalidForm
		}

		board = domain.Board{
			Name:   form.Name,
			Width:  form.Width,
			Height: form.Height,
		}

		if err := tx.SaveBoard(&board); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &board, nil
}

func Update(id interface{}, form *Form) (board *domain.Board, err error) {
	form.NormalizeInputs()

	err = repo.Transaction(func (tx Repository) (err error) {
		if !form.IsValid(tx) {
			return domain.ErrInvalidForm
		}

		board, err = tx.GetBoardByID(id)
		if err != nil {
			return err
		}

		board.Name = form.Name
		board.Width = form.Width
		board.Height = form.Height

		err = tx.SaveBoard(board)
		if err != nil {
			return err
		}

		return
	})

	return
}

func UpdateName(id interface{}, form *Form) (board *domain.Board, err error) {
	form.NormalizeInputs()

	err = repo.Transaction(func (tx Repository) (err error) {
		if !form.IsValid(tx) {
			return domain.ErrInvalidForm
		}

		board, err = tx.GetBoardByID(id)
		if err != nil {
			return err
		}

		board.Name = form.Name

		err = tx.SaveBoard(board)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func UpdateDimensions(id interface{}, form *Form) (board *domain.Board, err error) {
	form.NormalizeInputs()

	err = repo.Transaction(func (tx Repository) (err error) {
		board, err = repo.GetBoardByID(id)
		if err != nil {
			return err
		}

		form.Name = board.Name // ensure that a blank name does not trigger a validation error
		if !form.IsValid(tx) {
			return domain.ErrInvalidForm
		}

		board.Width = form.Width
		board.Height = form.Height

		err = tx.SaveBoard(board)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func DeleteByID(id interface{}) error {
	return repo.DeleteBoardByID(id)
}

