package board

import (
	"city-route-game/domain"
	"github.com/assertgo/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"testing"
)

var dbConn *gorm.DB

func TestMain(m *testing.M) {
	dbPath := "file::memory:?cache=shared"

	var err error
	dbConn, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
		DisableNestedTransaction: true,
	})
	if err != nil {
		panic("Error connecting to test db: " + err.Error())
	}

	if err = dbConn.AutoMigrate(domain.Models()...); err != nil {
		panic("Error migrating test db: " + err.Error())
	}

	Init(NewGormRepositoryImpl(dbConn))

	os.Exit(m.Run())
}

func TestCreateWithValidForm(t *testing.T) {
	assert := assert.New(t)

	form := Form{
		Name:   "My Awesome Board",
		Width:  200,
		Height: 300,
	}

	board, err := Create(&form)

	assert.That(err).IsNil()
	assert.That(board).IsNotNil()
	assert.ThatString(board.Name).IsEqualTo("My Awesome Board")
	assert.ThatInt(board.Width).IsEqualTo(200)
	assert.ThatInt(board.Height).IsEqualTo(300)
}

func TestUpdateWithValidForm(t *testing.T) {
	assert := assert.New(t)

	board, err := Create(&Form{
		Name:   "Original Name",
		Width:  200,
		Height: 300,
	})
	assert.That(err).IsNil()

	form := NewUpdateForm(board)
	form.Name = "New Name"
	form.Width = 213
	form.Height = 453

	_, err = Update(board.ID, &form)
	assert.That(err).IsNil()

	board, err = FindByID(board.ID)
	assert.That(err).IsNil()

	assert.ThatString(board.Name).IsEqualTo(form.Name)
	assert.ThatInt(board.Width).IsEqualTo(form.Width)
	assert.ThatInt(board.Height).IsEqualTo(form.Height)
}

func TestUpdateBoardWithInvalidForm(t *testing.T) {
	assert := assert.New(t)

	otherBoard, err := Create(&Form{
		Name: "Other Board",
	})
	assert.That(err).IsNil()

	originalBoard, err := Create(&Form{
		Name:   "Original Name",
		Width:  200,
		Height: 300,
	})
	assert.That(err).IsNil()

	form := Form{
		ID:     originalBoard.ID,
		Name:   otherBoard.Name,
		Width:  213,
		Height: 453,
	}

	_, err = Update(form.ID, &form)
	assert.That(err).IsNotNil()

	var updatedBoard *domain.Board
	updatedBoard, err = FindByID(originalBoard.ID)
	assert.That(err).IsNil()

	assert.ThatString(updatedBoard.Name).IsEqualTo(originalBoard.Name)
	assert.ThatInt(updatedBoard.Width).IsEqualTo(originalBoard.Width)
	assert.ThatInt(updatedBoard.Height).IsEqualTo(originalBoard.Height)
}

func TestDeleteByID(t *testing.T) {
	assert := assert.New(t)

	board, err := Create(&Form{
		Name: "Test Board",
	})
	assert.That(err).IsNil()

	err = DeleteByID(board.ID)
	assert.That(err).IsNil()

	board, err = FindByID(board.ID)
	assert.That(err).IsEqualTo(gorm.ErrRecordNotFound)
}
