package app

import (
	"fmt"
	"strings"
)

func NewCreateBoardForm() CreateBoardForm {
	return CreateBoardForm{
		Form: NewPostForm("/boards/"),
		Name: "",
	}
}

type UpdateBoardForm struct {
	Form   `json:"-"`
	ID     ID     `json:"id" schema:"id"`
	Name   string `json:"name" schema:"name"`
	Width  int    `json:"width" schema:"width"`
	Height int    `json:"height" schema:"height"`
}

func NewUpdateBoardForm(board *Board) UpdateBoardForm {
	return UpdateBoardForm{
		Form: Form{
			Errors: nil,
			Action: fmt.Sprintf("/boards/%d", board.ID),
			Method: "PATCH",
		},
		ID:     board.ID,
		Name:   board.Name,
		Width:  board.Width,
		Height: board.Height,
	}
}

func NewBoardNameForm(board *Board) UpdateBoardForm {
	return UpdateBoardForm{
		Form: Form{
			Action: fmt.Sprintf("/boards/%d", board.ID),
			Method: "PATCH",
		},
		ID:   board.ID,
		Name: board.Name,
	}
}

type CreateBoardForm struct {
	Form `json:"-"`
	ID ID `json:"id" schema:"ID"`
	Name string `json:"name" schema:"Name"`
}

// CityForm JSON format in which cities will be posted from the board editor on create or update.
// Cities are always valid as long as they relate to a board;
// we let the user do whatever they want with them.
type CityForm struct {
	ID       uint     `json:"id" schema:"id"`
	Name     string   `json:"name" schema:"name"`
	Position Position `json:"position" schema:"position"`
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
