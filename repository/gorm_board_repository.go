package repository

import (
	"city-route-game/domain"
	"errors"
	"gorm.io/gorm"
)

func NewGormBoardRepository(db *gorm.DB) domain.BoardRepository {
	return &gormBoardRepository{
		db: db,
	}
}

type gormBoardRepository struct {
	db *gorm.DB
}

func (p *gormBoardRepository) GetBoardByID(id interface{}) (*domain.Board, error) {
	var board domain.Board
	if err := p.db.First(&board, id).Error; err != nil {
		return nil, err
	}
	return &board, nil
}

func (p *gormBoardRepository) CreateBoard(board *domain.Board) error {
	if err := p.db.Save(board).Error; err != nil {
		return err
	}
	return nil
}

func (p *gormBoardRepository) UpdateBoard(id interface{}, updateFn func (board *domain.Board) (*domain.Board, error)) error {
	return p.db.Transaction(func (tx *gorm.DB) error {
		var board domain.Board
		err := tx.First(&board, id).Error
		if err != nil {
			return err
		}

		updatedBoard, err := updateFn(&board)
		if err != nil {
			return err
		}
		if updatedBoard == nil {
			return errors.New("updateFn returned nil board and nil error")
		}

		err = tx.Save(updatedBoard).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (p *gormBoardRepository) ListBoards() ([]domain.Board, error) {
	var boards []domain.Board
	if err := p.db.Find(&boards).Error; err != nil {
		return nil, err
	}
	return boards, nil
}

func (p *gormBoardRepository) DeleteBoardByID(id interface{}) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var board domain.Board
		var err error

		if err = p.db.First(&board, id).Error; err != nil {
			return err
		}

		if err = tx.Delete(&board).Error; err != nil {
			return err
		}

		return nil
	})
}

//func (p *gormBoardRepository) BoardExistsWithName(name string) (bool, error) {
//	var dupe domain.Board
//	err := p.db.Where("name = ?", name).Take(&dupe).Error
//	if err == nil {
//		return true, nil
//	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
//		return false, err
//	} else {
//		return false, nil
//	}
//}
//
//func (p *gormBoardRepository) BoardExistsWithNameAndIdNot(name string, idNot interface{}) (bool, error) {
//	var dupe domain.Board
//	err := p.db.Where("name = ? AND id <> ?", name, idNot).Take(&dupe).Error
//	if err == nil {
//		return true, nil
//	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
//		return false, err
//	} else {
//		return false, nil
//	}
//}

func (p *gormBoardRepository) ListCitiesByBoardID(boardID interface{}) ([]domain.City, error) {
	var cities []domain.City
	err := p.db.Preload("CitySpaces").
		Where("`cities`.`board_id` = ?", boardID).
		Find(&cities).
		Error
	if err != nil {
		return nil, err
	}

	//cityIDs := make([]uint, len(cities))
	//for i, city := range cities {
	//	cityIDs[i] = city.ID
	//}

	//var spaces []domain.CitySpace
	//err = p.db.Where("city_id in (?)", cityIDs).Find(&spaces).Order("city_id, order").Error
	//
	//for

	return cities, nil
}

func (p *gormBoardRepository) GetCityByID(id interface{}) (*domain.City, error) {
	var city domain.City

	err := p.db.Preload("CitySpaces", func(spaces *gorm.DB) *gorm.DB {
		return spaces.Order("`city_spaces`.`order` ASC")
	}).First(&city, id).Error
	if err != nil {
		return nil, err
	}

	return &city, nil
}

func (p *gormBoardRepository) CreateCity(city *domain.City) error {
	if err := p.db.Save(city).Error; err != nil {
		return err
	}
	return nil
}

func (p *gormBoardRepository) UpdateCity(id interface{}, updateFn func (city *domain.City) (*domain.City, error) ) error {
	return p.db.Transaction(func (tx *gorm.DB) error {
		var city domain.City
		err := tx.First(&city, id).Error
		if err != nil {
			return err
		}

		updatedCity, err := updateFn(&city)
		if err != nil {
			return err
		}
		if updatedCity == nil {
			return errors.New("updateFn returned nil error and nil city")
		}

		err = tx.Save(updatedCity).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (p *gormBoardRepository) DeleteCityByBoardIDAndCityID(boardID interface{}, cityID interface{}) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var city domain.City
		if err := tx.First(&city, "board_id = ? AND id = ?", boardID, cityID).Error; err != nil {
			return err
		}

		if err := tx.Delete(&domain.CitySpace{}, "city_id = ?", cityID).Error; err != nil {
			return err
		}

		return p.DeleteCity(&city)
	})
}

func (p *gormBoardRepository) DeleteCity(city *domain.City) error {
	return p.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Delete(city).Error; err != nil {
			return err
		}

		return nil
	})
}

func (p *gormBoardRepository) CreateCitySpace(citySpace *domain.CitySpace) error {
	return p.db.Save(citySpace).Error
}

func (p *gormBoardRepository) UpdateCitySpace(id interface{}, updateFn func (*domain.CitySpace) (*domain.CitySpace, error)) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var citySpace *domain.CitySpace

		err := tx.First(citySpace, id).Error
		if err != nil {
			return err
		}

		citySpace, err = updateFn(citySpace)
		if err != nil {
			return err
		}

		err = tx.Save(citySpace).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (p *gormBoardRepository) GetCitySpacesByCityID(cityID interface{}) ([]domain.CitySpace, error) {
	var spaces []domain.CitySpace
	if err := p.db.Find(&spaces, "city_id = ?", cityID).Error; err != nil {
		return nil, err
	}
	return spaces, nil
}

func (p *gormBoardRepository) DeleteCitySpaceByID(id interface{}) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var space domain.CitySpace
		var err error

		if err = tx.First(&space, id).Error; err != nil {
			return err
		}

		if err := tx.Delete(&space).Error; err != nil {
			return err
		}

		return nil
	})
}
