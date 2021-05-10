package admin

import (
	"city-route-game/domain"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var ErrInvalidForm = errors.New("invalid form error")

type Form struct {
	errors map[string][]string `schema:"-"`
	Action string              `schema:"-"`
	Method string              `schema:"_method"`
}

func (f *Form) IsUpdate() bool {
	return f.Method == "PATCH" || f.Method == "PUT"
}

func (f *Form) IsInsert() bool {
	return f.Method == "POST"
}

func (f *Form) IsCreate() bool {
	return f.IsInsert()
}

func (f *Form) Errors() map[string][]string {
	return f.errors
}

func (f *Form) AddError(name string, msg string) {
	if f.errors == nil {
		f.errors = make(map[string][]string)
	}

	_, exists := f.errors[name]
	if !exists {
		f.errors[name] = make([]string, 0, 1)
	}

	f.errors[name] = append(f.errors[name], msg)
}

func (f *Form) ClearErrors() {
	if f.HasError() {
		f.errors = make(map[string][]string)
	}
}

func (f *Form) HasError() bool {
	if f.errors == nil {
		return false
	}

	if len(f.errors) == 0 {
		return false
	}

	return true
}

type BoardForm struct {
	Form
	ID     uint
	Name   string
	Width  int
	Height int
}

func (f *BoardForm) NormalizeInputs() {
	f.Name = strings.TrimSpace(f.Name)
}

func (f *BoardForm) IsValid(db *gorm.DB) bool {
	f.ClearErrors()

	if len(f.Name) == 0 {
		f.AddError("Name", "must not be blank")
	} else if len(f.Name) > 100 {
		f.AddError("Name", "is too long; must be 100 characters or less")
	} else {
		// Name is in range, so go ahead and check for duplicates
		var dupe domain.Board

		var query string
		if f.ID == 0 {
			query = "name = ?"
		} else {
			query = "name = ? AND id <> ?"
		}

		conditions := []interface{}{
			f.Name,
		}
		if f.ID != 0 {
			conditions = append(conditions, f.ID)
		}

		err := db.Where(query, conditions...).First(&dupe).Error
		if err == nil {
			f.AddError("Name", "has already been taken")
		} else if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}

	return !f.HasError()
}

// JSON format in which cities will be posted from the board editor on create or update.
// Cities are always valid as long as they relate to a board;
// we let the user do whatever they want with them.
type CityForm struct {
	ID       uint            `json:"id"`
	Name     string          `json:"name"`
	Position domain.Position `json:"position"`
}
