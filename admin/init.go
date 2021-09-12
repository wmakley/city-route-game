package admin

import (
	"city-route-game/internal/app"
	"github.com/gorilla/schema"
)

var (
	formDecoder     *schema.Decoder
	boardRepository app.BoardCrudRepository
	config          Config
)

type Config struct {
	TemplateRoot string
	AssetHost    string
	IPWhitelist  []string
}

func Init(config_ Config, boardRepo app.BoardCrudRepository) {
	formDecoder = schema.NewDecoder()
	boardRepository = boardRepo
	config = config_
}
