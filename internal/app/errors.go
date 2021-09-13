package app

import (
	"errors"
	"fmt"
)

// ErrInvalidForm Error to be returned by BoardEditorService on invalid user input
var ErrInvalidForm = errors.New("invalid form (check form object errors)")

// ErrNameTaken Error to be returned by BoardCrudRepository upon attempt to create
// or update a board with a duplicate name
var ErrNameTaken = errors.New("name already taken")

type RecordNotFound struct {
	Name string
	ID ID
}

func (e RecordNotFound) Error() string {
	return fmt.Sprint(e.Name, " with id ", e.ID, " not found")
}

func (e RecordNotFound) Is(target error) bool {
	_, sameType := target.(*RecordNotFound)
	return sameType
}

func NewBoardNotFoundError(id ID) error {
	return &RecordNotFound{
		Name: "Board",
		ID: id,
	}
}

func NewRecordNotFoundError(name string, id ID) error {
	return &RecordNotFound{
		Name: name,
		ID: id,
	}
}


type ErrInvalidIDString struct {
	Msg string
	Cause error
}

func (e ErrInvalidIDString) Error() string {
	return e.Msg
}

func (e ErrInvalidIDString) Is(target error) bool {
	_, ok := target.(*ErrInvalidIDString)
	return ok
}

func (e ErrInvalidIDString) Unwrap() error {
	return e.Cause
}
