package domain

import (
	"city-route-game/repository"
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"testing"
)

var (
	Rollback = errors.New("rollback")
)

func TestMain(m *testing.M) {
	dbPath := "file::memory:?cache=shared"

	dbConn, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic("Error connecting to repository: " + err.Error())
	}

	if err := dbConn.AutoMigrate(Models()...); err != nil {
		panic("Error migrating repository: " + err.Error())
	}

	Init(Config{
		BoardRepository: repository.NewGormProvider(dbConn),
	})

	os.Exit(m.Run())
}

