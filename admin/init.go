package admin

import (
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB
	formDecoder *schema.Decoder
)

func Init(dbConn *gorm.DB) {
	db = dbConn
	formDecoder = schema.NewDecoder()
}
