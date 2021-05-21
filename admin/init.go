package admin

import (
	"github.com/gorilla/schema"
)

var (
	formDecoder *schema.Decoder
	config      Config
)

type Config struct {
	TemplateRoot string
	AssetHost    string
	IPWhitelist  []string
}

func Init(config_ Config) {
	formDecoder = schema.NewDecoder()
	config = config_
}
