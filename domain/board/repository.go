package board

import (
	"city-route-game/domain"
	"errors"
	"gorm.io/gorm"
)

type BoardRepository interface {
	GetBoardByID(id interface{}) (*domain.Board, error)
	SaveBoard(board *domain.Board) error
	ListBoardNamesAndIDs() ([]domain.Board, error)
	DeleteBoardByID(id interface{}) error
}

type GormRepositoryImpl struct {
	db *gorm.DB
}

func NewGormRepositoryImpl(dbConn *gorm.DB) *GormRepositoryImpl {
	return &GormRepositoryImpl{
		db: dbConn,
	}
}

func (p *GormRepositoryImpl)GetBoardByID(id interface{}) (*domain.Board, error) {
	var board domain.Board
	if err := p.db.First(&board, id).Error; err != nil {
		return nil, err
	}
	return &board, nil
}

func (p *GormRepositoryImpl)SaveBoard(board *domain.Board) error {
	if err := p.db.Save(board).Error; err != nil {
		return err
	}
	return nil
}

func (p *GormRepositoryImpl)ListBoards() ([]domain.Board, error) {
	var boards []domain.Board
	if err := p.db.Find(&boards).Error; err != nil {
		return nil, err
	}
	return boards, nil
}

func (p *GormRepositoryImpl)DeleteBoardByID(id interface{}) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var board domain.Board
		var err error

		if err = tx.First(&board, id).Error; err != nil {
			return err
		}

		//var cities []domain.City
		//if err = tx.Find(&cities, "board_id = ?", board.ID).Error; err != nil {
		//	return err
		//}
		//
		//for _, city := range cities {
		//	if err = tx.Delete(&city).Error; err != nil {
		//		return err
		//	}
		//}

		if err = tx.Delete(&board).Error; err != nil {
			return err
		}

		return nil
	})
}

func (p *GormRepositoryImpl)BoardExistsWithName(name string) (bool, error) {
	var dupe domain.Board
	err := p.db.Where("name = ?", name).Take(&dupe).Error
	if err == nil {
		return true, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	} else {
		return false, nil
	}
}

func (p *GormRepositoryImpl)BoardExistsWithNameAndIdNot(name string, idNot interface{}) (bool, error) {
	var dupe domain.Board
	err := p.db.Where("name = ? AND id <> ?", name, idNot).Take(&dupe).Error
	if err == nil {
		return true, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	} else {
		return false, nil
	}
}
