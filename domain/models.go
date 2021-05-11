package domain

import (
	"time"
)

// Return an empty instance of every model for use with gorm Automigration
func Models() []interface{} {
	return []interface{}{
		&Game{}, &Board{}, &Player{}, &PlayerBoard{}, &PlayerBonusToken{}, &BonusToken{}, &RouteBonusToken{}, &City{}, &CitySpace{}, &Route{}, &RouteSpace{},
	}
}

// Simpler version of gorm.Model with JSON tags and without the DeletedAt column.
// When we delete, we mean it!
type Model struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Shared mixin
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
}

// Board structure
type City struct {
	Model
	BoardID    uint   `json:"boardId" gorm:"not null;index"`
	Name       string `json:"name" gorm:"not null"`
	Position   `json:"position"`
	CitySpaces []CitySpace `json:"spaces" gorm:"foreignKey:CityID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Board structure
type CitySpace struct {
	Model
	CityID            uint          `json:"cityId" gorm:"not null;uniqueIndex:uidx_city_space_city_id_order"`
	Order             int           `json:"order" gorm:"not null;index:uidx_city_space_city_id_order"`
	SpaceType         TradesmanType `json:"spaceType" gorm:"not null;default:1"`
	RequiredPrivilege int           `json:"requiredPrivilege" gorm:"not null;default:1"`
}

// Board structure
type Route struct {
	Model
	StartCityID uint         `json:"startCityId" gorm:"not null;index"`
	EndCityID   uint         `json:"endCityId" gorm:"not null;index"`
	TavernFlag  bool         `json:"tavernFlag" gorm:"not null;default:0"`
	RouteSpaces []RouteSpace `json:"spaces"`
}

// Board structure
type RouteSpace struct {
	Model
	RouteID uint `json:"routeId" gorm:"not null;uniqueIndex:uidx_route_space_route_order"`
	Order   int  `json:"order" gorm:"not null;index:uidx_route_space_route_order"`
}

// Game state
type Game struct {
	Model
	Name             string `json:"name" gorm:"not null;index"`
	Coellen1PlayerID *uint
	Coellen2PlayerID *uint
	Coellen3PlayerID *uint
	Coellen4PlayerID *uint
}

// Game state
type Player struct {
	Model
	GameID uint   `json:"gameId" gorm:"uniqueIndex:uidx_game_id_player_name"`
	Name   string `json:"name" gorm:"not null;index:uidx_game_id_player_name"`
	Color  string `json:"color" gorm:"not null"`
	Score  int    `json:"score" gorm:"not null;default:0"`
}

// Game state
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
	GameID       uint `gorm:"not null;uniqueIndex:uidx_supply_bonus_token"`
	BonusTokenID uint `gorm:"not null;index:uidx_supply_bonus_token"`
	Order        int  `gorm:"not null;index:uidx_supply_bonus_token"`
	BonusToken   BonusToken
}

// Game state
type RouteBonusToken struct {
	Model
	GameID       uint `gorm:"not null;uniqueIndex:uidx_route_bonus_token"`
	RouteID      uint `gorm:"not null;index:uidx_route_bonus_token"`
	BonusTokenID uint `gorm:"not null;index:uidx_route_bonus_token"`
	BonusToken   BonusToken
}

// Game state
// Represents a single bonus token
type BonusToken struct {
	Model
	BonusTokenTypeID uint `gorm:"not null"`
	Gold             bool `gorm:"not null"`
}
