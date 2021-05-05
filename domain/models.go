package domain

import (
	"time"
)

// Return an empty instance of every model for use with gorm Automigration
func Models() []interface{} {
	return []interface{}{
		&Game{}, &Board{}, &Player{}, &PlayerBoard{}, &PlayerBonusToken{}, &BonusToken{}, &RouteBonusToken{}, &City{}, &CitySlot{}, &Route{}, &RouteSlot{},
	}
}

// Simpler version of gorm.Model with JSON tags and without the DeletedAt column.
// When we delete, we mean it!
type Model struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Game struct {
	Model
	Name             string `json:"name" gorm:"not null;index"`
	Coellen1PlayerID *uint
	Coellen2PlayerID *uint
	Coellen3PlayerID *uint
	Coellen4PlayerID *uint
}

type Player struct {
	Model
	GameID uint   `json:"gameId" gorm:"uniqueIndex:uidx_game_id_player_name"`
	Name   string `json:"name" gorm:"not null;index:uidx_game_id_player_name"`
	Color  string `json:"color" gorm:"not null"`
	Score  int    `json:"score" gorm:"not null;default:0"`
}

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
	PlateBonusToken   *BonusToken `gorm:"foreignKey:PlateBonusTokenID"`
}

// Join table between players and bonus tokens
type PlayerBonusToken struct {
	Model
	PlayerID     uint `json:"playerId" gorm:"uniqueIndex:uidx_player_bonus_token"`
	BonusTokenID uint `json:"bonusTokenId" gorm:"index:uidx_player_bonus_token"`
	BonusToken   BonusToken
	Played       bool `json:"played" gorm:"not null;default:0"`
}

// Represents a bonus token in the supply, initialized at start of game
type SupplyBonusToken struct {
	Model
	GameID       uint `gorm:"not null;uniqueIndex:uidx_supply_bonus_token"`
	BonusTokenID uint `gorm:"not null;index:uidx_supply_bonus_token"`
	Order        int  `gorm:"not null;index:uidx_supply_bonus_token"`
	BonusToken   BonusToken
}

type RouteBonusToken struct {
	Model
	GameID       uint `gorm:"not null;uniqueIndex:uidx_route_bonus_token"`
	RouteID      uint `gorm:"not null;index:uidx_route_bonus_token"`
	BonusTokenID uint `gorm:"not null;index:uidx_route_bonus_token"`
	BonusToken   BonusToken
}

// Represents a single bonus token
type BonusToken struct {
	Model
	BonusTokenTypeID uint `gorm:"not null"`
	Gold             bool `gorm:"not null"`
}

type Board struct {
	Model
	Name   string `json:"name" gorm:"not null;uniqueIndex"`
	GameID *uint  `json:"gameId" gorm:"index"`
	Width  int    `json:"width" gorm:"not null;default:0"`
	Height int    `json:"height" gorm:"not null;default:0"`
	Cities []City `json:"cities"`
}

type Position struct {
	X int `json:"x" gorm:"not null;default:0"`
	Y int `json:"y" gorm:"not null;default:0"`
}

type City struct {
	Model
	BoardID   uint   `gorm:"not null;index"`
	Name      string `gorm:"not null"`
	Position  `json:"pos"`
	CitySlots []CitySlot `json:"slots"`
}

type CitySlot struct {
	Model
	CityID            uint   `gorm:"not null;index"`
	Order             int    `gorm:"not null"`
	SlotType          string `gorm:"not null"`
	RequiredPrivilege int    `gorm:"not null;default:1"`
}

type Route struct {
	Model
	StartCityID uint `gorm:"not null;index"`
	EndCityID   uint `gorm:"not null;index"`
	Tavern      bool `gorm:"not null;default:0"`
	RouteSlots  []RouteSlot
}

type RouteSlot struct {
	Model
	RouteID uint `gorm:"not null;uniqueIndex:uidx_route_slot_route_order"`
	Order   int  `gorm:"not null;index:uidx_route_slot_route_order"`
}
