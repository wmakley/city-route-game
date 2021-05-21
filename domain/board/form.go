package board

import (
	"city-route-game/domain"
	"fmt"
	"strings"
)

func NewCreateForm() Form {
	return Form{
		Form: domain.NewPostForm("/boards/"),
		Name: "",
	}
}

func NewUpdateForm(board *domain.Board) Form {
	return Form{
		Form: domain.Form{
			Action: fmt.Sprintf("/boards/%d", board.ID),
			Method: "PATCH",
		},
		ID:   board.ID,
		Name: board.Name,
	}
}

type Form struct {
	domain.Form   `json:"-"`
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func (f *Form) NormalizeInputs() {
	f.Name = strings.TrimSpace(f.Name)

	if f.Width <= 0 {
		f.Width = 800
	}

	if f.Height <= 0 {
		f.Height = 500
	}
}

func (f *Form) IsValid(repo Repository) bool {
	f.ClearErrors()

	if len(f.Name) == 0 {
		f.AddError("Name", "must not be blank")
	} else if len(f.Name) > 100 {
		f.AddError("Name", "is too long; must be 100 characters or less")
	} else {
		// Name is in range, so go ahead and check for duplicates
		var duped bool
		var err error
		if f.ID == 0 {
			duped, err = repo.BoardExistsWithName(f.Name)
		} else {
			duped, err = repo.BoardExistsWithNameAndIdNot(f.Name, f.ID)
		}

		if err != nil {
			panic(err)
		}

		if duped {
			f.AddError("Name", "has already been taken")
		}
	}

	if f.Width < 0 {
		f.AddError("Width", "cannot be less than 0")
	}
	if f.Height < 0 {
		f.AddError("Height", "cannot be less than 0")
	}

	return !f.HasError()
}
