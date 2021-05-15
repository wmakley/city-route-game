package admin

import (
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

var (
	db           *gorm.DB
	formDecoder  *schema.Decoder
	templateRoot string
	assetHost    string
)

func Init(dbConn *gorm.DB, templateRoot_ string, assetHost_ string) {
	db = dbConn
	formDecoder = schema.NewDecoder()
	templateRoot = templateRoot_
	assetHost = assetHost_
}
