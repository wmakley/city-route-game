package domain

import "gorm.io/gorm"

var (
	DB *gorm.DB
)

func Init(DB_ *gorm.DB) {
	DB = DB_
}
