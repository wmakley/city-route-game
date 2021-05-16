package admin

import (
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB
	formDecoder *schema.Decoder
	config      Config
)

type Config struct {
	TemplateRoot string
	AssetHost    string
	IPWhitelist  []string
}

func Init(dbConn *gorm.DB, templateRoot string, assetHost string, ipWhitelist []string) {
	db = dbConn
	formDecoder = schema.NewDecoder()

	config = Config{
		TemplateRoot: templateRoot,
		AssetHost:    assetHost,
		IPWhitelist:  ipWhitelist,
	}
}
