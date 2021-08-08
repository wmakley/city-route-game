package city

import (
	"city-route-game/domain"
	"strings"

	"gorm.io/gorm"
)

// Form JSON format in which cities will be posted from the board editor on create or update.
// Cities are always valid as long as they relate to a board;
// we let the user do whatever they want with them.
type Form struct {
	ID       uint            `json:"id"`
	Name     string          `json:"name"`
	Position domain.Position `json:"position"`
}

func (f *Form) NormalizeInputs() {
	f.Name = strings.TrimSpace(f.Name)
}

func (f *Form) IsValid(tx *gorm.DB) bool {
	return true
}

type AddCitySpaceForm struct {
	CityID            uint
	SpaceType         domain.TradesmanType
	RequiredPrivilege int
	domain.Form
}

func (f *AddCitySpaceForm) IsValid(tx *gorm.DB) bool {
	if f.CityID == 0 {
		f.AddError("CityID", "is required")
	}

	if f.RequiredPrivilege < 1 || f.RequiredPrivilege > 4 {
		f.AddError("RequiredPrivilege", "is out of bounds (must be between 1 and 4)")
	}

	if f.SpaceType != domain.TraderID && f.SpaceType != domain.MerchantID {
		f.AddError("SpaceType", "is invalid")
	}

	return !f.HasError()
}

type UpdateCitySpaceForm struct {
	ID                uint
	SpaceType         domain.TradesmanType
	RequiredPrivilege int
	domain.Form
}

func (f *UpdateCitySpaceForm) IsValid(tx *gorm.DB) bool {
	if f.RequiredPrivilege < 1 || f.RequiredPrivilege > 4 {
		f.AddError("RequiredPrivilege", "is out of bounds (must be between 1 and 4)")
	}

	if f.SpaceType != domain.TraderID && f.SpaceType != domain.MerchantID {
		f.AddError("SpaceType", "is invalid")
	}

	return f.HasError()
}
