package admin

import (
	"city-route-game/domain"
	"github.com/gorilla/schema"
)

var (
	formDecoder *schema.Decoder
	boardRepository *domain.BoardRepository
	config      Config
)

type Config struct {
	TemplateRoot string
	AssetHost    string
	IPWhitelist  []string
}

func Init(config_ Config, boardRepo *domain.BoardRepository) {
	formDecoder = schema.NewDecoder()
	boardRepository = boardRepo
	config = config_
}
