package gorm_board_crud_repository

import (
	"city-route-game/internal/app"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// Models Return an empty instance of every model for use with gorm Automigration
func models() []interface{} {
	return []interface{}{
		&Game{},
		&Board{},
		&Player{},
		&PlayerBoard{},
		&PlayerBonusToken{},
		&BonusToken{},
		&RouteBonusToken{},
		&City{},
		&CitySpace{},
		&Route{},
		&RouteSpace{},
	}
}

type constraintViolation struct {
	msg string
}

func (c *constraintViolation)Error() string {
	return c.msg
}

// Model is a simpler version of gorm.Model with JSON tags and without the DeletedAt column.
// When we delete, we mean it!
type Model struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Position is a shared mixin with X and Y
type Position struct {
	X int `json:"x" gorm:"not null;default:0"`
	Y int `json:"y" gorm:"not null;default:0"`
}

// Board structure base model
type Board struct {
	Model
	Name   string `json:"name" gorm:"not null;uniqueIndex"`
	GameID *uint  `json:"gameId" gorm:"index"`
	Width  int    `json:"width" gorm:"not null;default:0"`
	Height int    `json:"height" gorm:"not null;default:0"`
	Cities []City `json:"cities"`
}

func newGormBoardFromDomainBoard(board *app.Board) (*Board, error) {
	if board == nil {
		panic("board must not be nil")
	}

	return &Board{
		Model: Model{
			ID:        board.ID,
			CreatedAt: board.CreatedAt,
			UpdatedAt: board.UpdatedAt,
		},
		Name: board.Name,
		Width: board.Width,
		Height: board.Height,
	}, nil
}

func newDomainBoardFromGormBoard(gormBoard *Board) *app.Board {
	if gormBoard == nil {
		panic("gormBoard must not be nil")
	}

	return &app.Board{
		Model: app.Model{
			ID:        gormBoard.ID,
			CreatedAt: gormBoard.CreatedAt,
			UpdatedAt: gormBoard.UpdatedAt,
		},
		Name: gormBoard.Name,
		Width: gormBoard.Width,
		Height: gormBoard.Height,
	}
}

func (b *Board)BeforeDelete(tx *gorm.DB) error {
	var cities []City
	var err error

	err = tx.Find(&cities, "board_id = ?", b.ID).Error
	if err != nil {
		return err
	}

	for _, city := range cities {
		if err = tx.Delete(&city).Error; err != nil {
			return err
		}
	}

	return nil
}

// City part of the Board structure
type City struct {
	Model
	BoardID    uint   `json:"boardId" gorm:"not null;index"`
	Name       string `json:"name" gorm:"not null"`
	Position   `json:"position"`
	CitySpaces []CitySpace `json:"spaces"`
}

func newGormCityFromAppCity(appCity *app.City)(*City, error) {
	if appCity == nil {
		panic("appCity must not be nil")
	}

	city := City{
		Model: Model{
			ID: appCity.ID,
			CreatedAt: appCity.CreatedAt,
			UpdatedAt: appCity.UpdatedAt,
		},
		BoardID: appCity.BoardID,
		Name: appCity.Name,
		Position: Position{
			X: appCity.Position.X,
			Y: appCity.Position.Y,
		},
		CitySpaces: nil,
	}

	if appCity.CitySpaces != nil {
		city.CitySpaces = make([]CitySpace, 0, len(appCity.CitySpaces))
		for _, space := range appCity.CitySpaces {
			gormSpace, err := newGormCitySpaceFromAppCitySpace(&space)
			if err != nil {
				return nil, err
			}
			city.CitySpaces = append(city.CitySpaces, *gormSpace)
		}
	}

	return &city, nil
}

func newAppCityFromGormCity(gormCity *City) *app.City {
	if gormCity == nil {
		panic("gormCity must not be nil")
	}

	city := app.City{
		Model: app.Model{
			ID: gormCity.ID,
			CreatedAt: gormCity.CreatedAt,
			UpdatedAt: gormCity.UpdatedAt,
		},
		BoardID: gormCity.BoardID,
		Name: gormCity.Name,
		Position: app.Position{
			X: gormCity.Position.X,
			Y: gormCity.Position.Y,
		},
		CitySpaces: nil,
	}

	if gormCity.CitySpaces != nil {
		city.CitySpaces = make([]app.CitySpace, 0, len(gormCity.CitySpaces))
		for _, space := range gormCity.CitySpaces {
			appSpace := *newAppCitySpaceFromGormCitySpace(&space)
			city.CitySpaces = append(city.CitySpaces, appSpace)
		}
	}

	return &city
}

func (c *City)BeforeSave(tx *gorm.DB) error {
	// Ensure board exists
	var result []uint
	err := tx.Table("boards").
		Where("id = ?", c.BoardID).
		Limit(1).
		Pluck("id", &result).Error
	if err != nil {
		return err
	}
	if len(result) <= 0 {
		return fmt.Errorf("constraint violation: board with id %d not found", c.BoardID)
	}
	return nil
}

func (c *City)BeforeDelete(tx *gorm.DB) error {
	err := tx.Delete(&CitySpace{}, "city_id = ?", c.ID).Error
	if err != nil {
		return err
	}

	var routes []Route
	if err = tx.Find(&routes, "start_city_id = ? OR end_city_id = ?", c.ID, c.ID).Error; err != nil {
		return err
	}

	for _, route := range routes {
		if err = tx.Delete(&route).Error; err != nil {
			return err
		}
	}

	return nil
}

// CitySpace Part of a City, which is part of Board
type CitySpace struct {
	Model
	CityID            uint          `json:"cityId" gorm:"not null;uniqueIndex:uidx_city_space_city_id_order"`
	Order             int           `json:"order" gorm:"not null;index:uidx_city_space_city_id_order"`
	SpaceType         app.TradesmanType `json:"spaceType" gorm:"not null;default:1"`
	RequiredPrivilege int           `json:"requiredPrivilege" gorm:"not null;default:1"`
}

func newGormCitySpaceFromAppCitySpace(space *app.CitySpace)(*CitySpace, error) {
	return &CitySpace{
		Model: Model{
			ID: space.ID,
			CreatedAt: space.CreatedAt,
			UpdatedAt: space.UpdatedAt,
		},
		CityID: space.CityID,
		SpaceType: space.SpaceType,
		RequiredPrivilege: space.RequiredPrivilege,
		Order: space.Order,
	}, nil
}

func newAppCitySpaceFromGormCitySpace(space *CitySpace)(*app.CitySpace) {
	return &app.CitySpace{
		Model: app.Model{
			ID: space.ID,
			CreatedAt: space.CreatedAt,
			UpdatedAt: space.UpdatedAt,
		},
		CityID: space.CityID,
		SpaceType: space.SpaceType,
		RequiredPrivilege: space.RequiredPrivilege,
		Order: space.Order,
	}
}

// Route Connects two City on a Board
type Route struct {
	Model
	StartCityID uint         `json:"startCityId" gorm:"not null;index"`
	EndCityID   uint         `json:"endCityId" gorm:"not null;index"`
	TavernFlag  bool         `json:"tavernFlag" gorm:"not null;default:0"`
	RouteSpaces []RouteSpace `json:"spaces"`
}

func (r *Route)BeforeDelete(tx *gorm.DB) error {
	if err := tx.Delete(&RouteSpace{}, "route_id = ?", r.ID).Error; err != nil {
		return err
	}
	return nil
}

// RouteSpace part of the board structure
type RouteSpace struct {
	Model
	RouteID uint `json:"routeId" gorm:"not null;uniqueIndex:uidx_route_space_route_order"`
	Order   int  `json:"order" gorm:"not null;index:uidx_route_space_route_order"`
}

// Game represents the game state
type Game struct {
	Model
	Name             string `json:"name" gorm:"not null;index"`
	Coellen1PlayerID *uint
	Coellen2PlayerID *uint
	Coellen3PlayerID *uint
	Coellen4PlayerID *uint
}

// Player is part of the game state
type Player struct {
	Model
	GameID uint   `json:"gameId" gorm:"uniqueIndex:uidx_game_id_player_name"`
	Name   string `json:"name" gorm:"not null;index:uidx_game_id_player_name"`
	Color  string `json:"color" gorm:"not null"`
	Score  int    `json:"score" gorm:"not null;default:0"`
}

// PlayerBoard part of the game state
// todo: unique index on game id and player id
type PlayerBoard struct {
	Model
	GameID            uint `json:"gameId" gorm:"uniqueIndex:uidx_player_board"`
	PlayerID          uint `json:"playerId" gorm:"index:uidx_player_board"`
	Merchants         int  `gorm:"not null;default:0"`
	Traders           int  `gorm:"not null;default:0"`
	MerchantSupply    int  `gorm:"not null;default:0"`
	TraderSupply      int  `gorm:"not null;default:0"`
	ActionLevel       int  `gorm:"not null;default:1"`
	BankLevel         int  `gorm:"not null;default:1"`
	MoveLevel         int  `gorm:"not null;default:1"`
	KnowledgeLevel    int  `gorm:"not null;default:1"`
	CityKeyLevel      int  `gorm:"not null;default:1"`
	PrivilegeLevel    int  `gorm:"not null;default:1"`
	PlateBonusTokenID *uint
}

// Game state
// Join table between players and bonus tokens
type PlayerBonusToken struct {
	Model
	PlayerID     uint `json:"playerId" gorm:"uniqueIndex:uidx_player_bonus_token"`
	BonusTokenID uint `json:"bonusTokenId" gorm:"index:uidx_player_bonus_token"`
	BonusToken   BonusToken
	Played       bool `json:"played" gorm:"not null;default:0"`
}

// Game state
// Represents a bonus token in the supply, initialized at start of game
type SupplyBonusToken struct {
	Model
	GameID     uint `gorm:"not null;uniqueIndex:uidx_supply_bonus_token"`
	BonusTokenID uint `gorm:"not null;index:uidx_supply_bonus_token"`
	Order      int  `gorm:"not null;index:uidx_supply_bonus_token"`
	BonusToken BonusToken
}

// Game state
type RouteBonusToken struct {
	Model
	GameID       uint `gorm:"not null;uniqueIndex:uidx_route_bonus_token"`
	RouteID      uint `gorm:"not null;index:uidx_route_bonus_token"`
	BonusTokenID uint `gorm:"not null;index:uidx_route_bonus_token"`
	BonusToken   BonusToken
}

// BonusToken represents a single bonus token in the game state
type BonusToken struct {
	Model
	BonusTokenTypeID uint `gorm:"not null"`
	Gold             bool `gorm:"not null"`
}

