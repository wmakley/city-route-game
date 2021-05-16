package domain

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

var ErrInvalidForm = errors.New("invalid form error")

func NewPostForm(action string) Form {
	return Form{
		Action: action,
		Method: "POST",
	}
}

type Form struct {
	errors map[string][]string `schema:"-" json:"-"`
	Action string              `schema:"-" json:"-"`
	Method string              `schema:"_method" json:"-"`
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

func NewBoardForm() BoardForm {
	return BoardForm{
		Form: NewPostForm("/boards/"),
		Name: "",
	}
}

func NewEditBoardForm(board *Board) BoardForm {
	return BoardForm{
		Form: Form{
			Action: fmt.Sprintf("/boards/%d", board.ID),
			Method: "PATCH",
		},
		ID:   board.ID,
		Name: board.Name,
	}
}

type BoardForm struct {
	Form   `json:"-"`
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
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
		var dupe Board

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
	ID       uint     `json:"id"`
	Name     string   `json:"name"`
	Position Position `json:"position"`
}

func (f *CityForm) NormalizeInputs() {
	f.Name = strings.TrimSpace(f.Name)
}

func (f *CityForm) IsValid() bool {
	return true
}

type CitySpaceForm struct {
	CityID            uint
	SpaceType         TradesmanType
	RequiredPrivilege int
	Form
}

func (f *CitySpaceForm) NormalizeInputs() {
}

func (f *CitySpaceForm) IsValid(tx *gorm.DB) bool {
	if f.CityID == 0 {
		f.AddError("CityID", "is required")
	}

	if f.RequiredPrivilege < 1 || f.RequiredPrivilege > 4 {
		f.AddError("RequiredPrivilege", "is out of bounds (must be between 1 and 4)")
	}

	if f.SpaceType != TraderID && f.SpaceType != MerchantID {
		f.AddError("SpaceType", "is invalid")
	}

	return f.HasError()
}
