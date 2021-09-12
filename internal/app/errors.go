package app

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

var ErrInvalidForm = errors.New("invalid form error")
var ErrNameTaken = errors.New("name already taken")

type RecordNotFound struct {
	name string
	id ID
	cause error
}

func (e RecordNotFound) Error() string {
	return fmt.Sprint(e.name, " with id ", e.id, " not found")
}

func (e RecordNotFound) Is(target error) bool {
	_, isDirectMatch := target.(RecordNotFound)
	if isDirectMatch {
		return true
	}

	// Treat gorm not found errors as though they are this error
	return errors.Is(gorm.ErrRecordNotFound, target)
}

func (e RecordNotFound) Unwrap() error {
	return e.cause
}

func NewBoardNotFoundError(id ID, cause error) error {
	return &RecordNotFound{
		name: "Board",
		id: id,
		cause: cause,
	}
}

func NewRecordNotFoundError(name string, id ID, cause error) error {
	return &RecordNotFound{
		name: name,
		id: id,
		cause: cause,
	}
}
