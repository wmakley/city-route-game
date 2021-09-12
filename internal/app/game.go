package app

// Player is part of the game state
type Player struct {
	Model
	GameID ID     `json:"gameId"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	Score  int    `json:"score"`
}

// PlayerBoard part of the game state
// todo: unique index on game id and player id
type PlayerBoard struct {
	Model
	GameID            ID  `json:"gameId"`
	PlayerID          ID  `json:"playerId"`
	Merchants         int `json:"merchants"`
	Traders           int  `json:"traders"`
	MerchantSupply    int  `json:"merchantSupply"`
	TraderSupply      int  `json:"traderSupply"`
	ActionLevel       int  `json:"actionLevel"`
	BankLevel         int  `json:"bankLevel"`
	MoveLevel         int  `json:"moveLevel"`
	KnowledgeLevel    int  `json:"knowledgeLevel"`
	CityKeyLevel      int  `json:"cityKeyLevel"`
	PrivilegeLevel    int `json:"privilegeLevel"`
	PlateBonusTokenID *ID `json:"plateBonusTokenID"`
}


// Game represents the game state
type Game struct {
	Model
	Name             string `json:"name"`
	Coellen1PlayerID *ID    `json:"coellen1PlayerID"`
	Coellen2PlayerID *ID    `json:"coellen2PlayerID"`
	Coellen3PlayerID *ID    `json:"coellen3PlayerID"`
	Coellen4PlayerID *ID    `json:"coellen4PlayerID"`
}

// Game state
// Join table between players and bonus tokens
type PlayerBonusToken struct {
	Model
	PlayerID     ID         `json:"playerId"`
	BonusTokenID ID         `json:"bonusTokenId" `
	BonusToken   BonusToken `json:"bonusToken"`
	Played       bool       `json:"played"`
}

// Game state
// Represents a bonus token in the supply, initialized at start of game
type SupplyBonusToken struct {
	Model
	GameID       ID         `json:"gameId"`
	BonusTokenID ID         `json:"bonusTokenID"`
	Order        int        `json:"order"`
	BonusToken   BonusToken `json:"bonusToken"`
}

// Game state
type RouteBonusToken struct {
	Model
	GameID       ID         `json:"gameId"`
	RouteID      ID         `json:"routeId"`
	BonusTokenID ID         `json:"bonusTokenID"`
	BonusToken   BonusToken `json:"bonusToken"`
}

// BonusToken represents a single bonus token in the game state
type BonusToken struct {
	Model
	BonusTokenTypeID ID   `json:"bonusTokenTypeID"`
	Gold             bool `json:"gold"`
}
