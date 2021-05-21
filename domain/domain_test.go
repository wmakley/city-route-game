package domain

import (
	"city-route-game/gorm_provider"
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
		panic("Error connecting to gorm_provider: " + err.Error())
	}

	if err := dbConn.AutoMigrate(Models()...); err != nil {
		panic("Error migrating gorm_provider: " + err.Error())
	}

	Init(Config{
		PersistenceProvider: gorm_provider.NewGormProvider(dbConn),
	})

	os.Exit(m.Run())
}

