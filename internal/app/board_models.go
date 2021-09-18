package app

import (
	"fmt"
	"strconv"
	"time"
)

type ID = uint64

func NewIDFromString(str string) (ID, error) {
	id, err := strconv.ParseUint(str, 0, 64)
	if err != nil {
		return 0, ErrInvalidIDString{
			Msg:   fmt.Sprintf("invalid ID string: %+s", err.Error()),
			Cause: err,
		}
	}
	return ID(id), nil
}

type TradesmanType = uint

const (
	TraderID TradesmanType = 1
	MerchantID TradesmanType = 2
)

// Board structure base model
type Board struct {
	Model
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Cities []City `json:"cities"`
}

// Model is a simpler version of gorm.Model with JSON tags and without the DeletedAt column.
// When we delete, we mean it!
type Model struct {
	ID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Position is a shared mixin with X and Y
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// City part of the Board structure
type City struct {
	Model
	BoardID    ID     `json:"boardId"`
	Name       string `json:"name"`
	Position   `json:"position"`
	CitySpaces []CitySpace `json:"spaces"`
}

// CitySpace Part of a City, which is part of Board
type CitySpace struct {
	Model
	CityID            ID                   `json:"cityId"`
	Order             int           `json:"order"`
	SpaceType         TradesmanType `json:"spaceType"`
	RequiredPrivilege int            `json:"requiredPrivilege"`
}

// Route Connects two City on a Board
type Route struct {
	Model
	StartCityID ID           `json:"startCityId"`
	EndCityID   ID           `json:"endCityId"`
	TavernFlag  bool         `json:"tavernFlag"`
	RouteSpaces []RouteSpace `json:"spaces"`
}

// RouteSpace part of the board structure
type RouteSpace struct {
	Model
	RouteID ID  `json:"routeId"`
	Order   int `json:"order"`
}
