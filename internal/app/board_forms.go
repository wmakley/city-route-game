package app

import (
	"fmt"
	"strings"
)



func NewCreateBoardForm() BoardNameForm {
	return BoardNameForm{
		Form: NewPostForm("/boards/"),
		Name: "",
	}
}

func NewBoardNameForm(board *Board) BoardNameForm {
	return BoardNameForm{
		Form: Form{
			Action: fmt.Sprintf("/boards/%d", board.ID),
			Method: "PATCH",
		},
		ID:   board.ID,
		Name: board.Name,
	}
}

type BoardNameForm struct {
	Form   `json:"-"`
	ID     uint   `json:"id"`
	Name   string `json:"name"`
}

func (f *BoardNameForm) NormalizeInputs() {
	f.Name = strings.TrimSpace(f.Name)
}

func (f *BoardNameForm) IsValid() bool {
	f.ClearErrors()

	if len(f.Name) == 0 {
		f.AddError("Name", "must not be blank")
	} else if len(f.Name) > 100 {
		f.AddError("Name", "is too long; must be 100 characters or less")
	}

	return !f.HasError()
}

func NewBoardDimensionsForm(board *Board) BoardDimensionsForm {
	return BoardDimensionsForm{
		Form: Form{
			Action: fmt.Sprintf("/boards/%d", board.ID),
			Method: "PATCH",
		},
		ID:   board.ID,
		Width: board.Width,
		Height: board.Height,
	}
}

type BoardDimensionsForm struct {
	Form `json:"-"`
	ID uint `json:"id"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// CityForm JSON format in which cities will be posted from the board editor on create or update.
// Cities are always valid as long as they relate to a board;
// we let the user do whatever they want with them.
type CityForm struct {
	ID       uint            `json:"id"`
	Name     string          `json:"name"`
	Position Position `json:"position"`
}

func (f *CityForm) NormalizeInputs() {
	f.Name = strings.TrimSpace(f.Name)
}

func (f *CityForm) IsValid() bool {
	return true
}

type AddCitySpaceForm struct {
	CityID            uint
	SpaceType         TradesmanType
	RequiredPrivilege int
	Form
}

func (f *AddCitySpaceForm) IsValid() bool {
	if f.CityID == 0 {
		f.AddError("CityID", "is required")
	}

	if f.RequiredPrivilege < 1 || f.RequiredPrivilege > 4 {
		f.AddError("RequiredPrivilege", "is out of bounds (must be between 1 and 4)")
	}

	if f.SpaceType != TraderID && f.SpaceType != MerchantID {
		f.AddError("SpaceType", "is invalid")
	}

	return !f.HasError()
}

type UpdateCitySpaceForm struct {
	ID                uint
	SpaceType         TradesmanType
	RequiredPrivilege int
	Form
}

func (f *UpdateCitySpaceForm) IsValid() bool {
	if f.RequiredPrivilege < 1 || f.RequiredPrivilege > 4 {
		f.AddError("RequiredPrivilege", "is out of bounds (must be between 1 and 4)")
	}

	if f.SpaceType != TraderID && f.SpaceType != MerchantID {
		f.AddError("SpaceType", "is invalid")
	}

	return f.HasError()
}
