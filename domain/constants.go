package domain

type ColorType string

type TradesmanType uint

type BonusTokenLocationName string

type BonusTokenType uint

const (
	// Player colors:
	ColorGreen  ColorType = "Green"
	ColorRed    ColorType = "Red"
	ColorBlue   ColorType = "Blue"
	ColorYellow ColorType = "Yellow"
	ColorPurple ColorType = "Purple"

	// Tradesman types:
	TraderID TradesmanType = iota + 1
	MerchantID

	// Bonus token types:
	BonusTokenExtraTradingPostID BonusTokenType = iota + 1
	BonusTokenExchangeTradingPostsID
	BonusTokenMove3TradesmenID
	BonusTokenDevelop1AbilityID
	BonusTokenPlusThreeActionsID
	BonusTokenPlusFourActionsID
)

var (
	BonusTokenTypeCounts map[BonusTokenType]int = map[BonusTokenType]int{
		BonusTokenExtraTradingPostID:     4,
		BonusTokenExchangeTradingPostsID: 3,
		BonusTokenMove3TradesmenID:       2,
		BonusTokenDevelop1AbilityID:      2,
		BonusTokenPlusThreeActionsID:     2,
		BonusTokenPlusFourActionsID:      2,
	}

	ActionsPerActionLevel map[int]int = map[int]int{
		1: 2,
		2: 3,
		3: 3,
		4: 4,
		5: 4,
		6: 5,
	}

	IncomePerBankLevel map[int]int = map[int]int{
		1: 3,
		2: 5,
		3: 7,
		4: 999,
	}

	MovesPerKnowledgeLevel map[int]int = map[int]int{
		1: 2,
		2: 3,
		3: 4,
		4: 5,
	}

	MultiplierPerCityKeysLevel map[int]int = map[int]int{
		1: 1,
		2: 2,
		3: 2,
		4: 3,
		5: 4,
	}
)
