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
	IPWhitelist  map[string]bool
}

func Init(dbConn *gorm.DB, templateRoot string, assetHost string, ipWhitelist []string) {
	db = dbConn
	formDecoder = schema.NewDecoder()

	ipWhitelistMap := make(map[string]bool, len(ipWhitelist))
	for _, ip := range ipWhitelist {
		ipWhitelistMap[ip] = true
	}

	config = Config{
		TemplateRoot: templateRoot,
		AssetHost:    assetHost,
		IPWhitelist:  ipWhitelistMap,
	}
}
