package gorm_provider

import (
	"city-route-game/domain"
	"errors"
	"gorm.io/gorm"
)

func NewGormProvider(db *gorm.DB) *GormProvider {
	return &GormProvider{
		DB: db,
		InTransaction: false,
	}
}

type GormProvider struct {
	DB *gorm.DB
	InTransaction bool
}

// Transaction start a new transaction, or reuse the current one
func (p *GormProvider)Transaction(callback func(domain.PersistenceProvider) error) error {
	if p.InTransaction {
		return callback(p)
	} else {
		return p.DB.Transaction(func (tx *gorm.DB) error {
			wrapper := GormProvider{
				DB: tx,
				InTransaction: true,
			}

			return callback(&wrapper)
		})
	}
}

func (p *GormProvider)GetBoardByID(id interface{}) (*domain.Board, error) {
	var board domain.Board
	if err := p.DB.First(&board, id).Error; err != nil {
		return nil, err
	}
	return &board, nil
}

func (p *GormProvider)SaveBoard(board *domain.Board) error {
	if err := p.DB.Save(board).Error; err != nil {
		return err
	}
	return nil
}

func (p *GormProvider)ListBoards() ([]domain.Board, error) {
	var boards []domain.Board
	if err := p.DB.Find(&boards).Error; err != nil {
		return nil, err
	}
	return boards, nil
}

func (p *GormProvider)DeleteBoardByID(id interface{}) error {
	return p.DB.Transaction(func(tx *gorm.DB) error {
		var board domain.Board
		var err error

		if err = p.DB.First(&board, id).Error; err != nil {
			return err
		}

		var cities []domain.City
		if err = p.DB.Find(&cities, "board_id = ?", board.ID).Error; err != nil {
			return err
		}

		for _, city := range cities {
			if err = p.DeleteCity(&city); err != nil {
				return err
			}
		}

		if err = tx.Delete(&board).Error; err != nil {
			return err
		}

		return nil
	})
}

func (p *GormProvider)BoardExistsWithName(name string) (bool, error) {
	var dupe domain.Board
	err := p.DB.Where("name = ?", name).Take(&dupe).Error
	if err == nil {
		return true, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	} else {
		return false, nil
	}
}

func (p *GormProvider)BoardExistsWithNameAndIdNot(name string, idNot interface{}) (bool, error) {
	var dupe domain.Board
	err := p.DB.Where("name = ? AND id <> ?", name, idNot).Take(&dupe).Error
	if err == nil {
		return true, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	} else {
		return false, nil
	}
}



func (p *GormProvider)SaveCitySpace(citySpace *domain.CitySpace) error {
	return p.DB.Transaction(func (tx *gorm.DB) error {
		if err := tx.Save(citySpace).Error; err != nil {
			return err
		}

		return nil
	})
}

func (p *GormProvider)GetCitySpacesByCityID(cityID interface{}) ([]domain.CitySpace, error) {
	var spaces []domain.CitySpace
	if err := p.DB.Find(&spaces, "city_id = ?", cityID).Error; err != nil {
		return nil, err
	}
	return spaces, nil
}
