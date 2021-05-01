package gamesapi

import "gorm.io/gorm"

var (
	db *gorm.DB
)

func Init(dbConn *gorm.DB) {
	db = dbConn
}
