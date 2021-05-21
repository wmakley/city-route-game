package city

import (
	"city-route-game/domain"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(dbConn *gorm.DB) {
	db = dbConn
}

func FindAllByBoardID(boardID interface{}) ([]domain.City, error) {
	var cities []domain.City
	err := db.Find(&cities, "board_id = ?", boardID).Error
	if err != nil {
		return nil, err
	}
	return cities, nil
}

func Create(boardID interface{}, form *Form) (*domain.City, error) {
	var city domain.City
	err := db.Transaction(func (tx *gorm.DB) error {
		form.NormalizeInputs()

		if !form.IsValid(tx) {
			return domain.ErrInvalidForm
		}

		var board domain.Board
		err := db.First(&board, boardID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("board id %v not found", boardID)
			} else {
				return err
			}
		}

		city = domain.City{
			BoardID:  board.ID,
			Name:     form.Name,
			Position: form.Position,
		}

		if err = tx.Save(&city).Error; err != nil {
			return err
		}

		return nil
	})
	return &city, err
}

func Update(form *Form) (city *domain.City, err error) {
	form.NormalizeInputs()

	err = db.Transaction(func(tx *gorm.DB) (err error) {
		if !form.IsValid(tx) {
			err = domain.ErrInvalidForm
			return
		}

		if err = tx.First(city, form.ID).Error; err != nil {
			return
		}

		city.Name = form.Name
		city.Position.X = form.Position.X
		city.Position.Y = form.Position.Y

		if err = tx.Save(city).Error; err != nil {
			return
		}

		return
	})

	return
}

func DeleteByBoardIDAndCityID(boardID interface{}, cityID interface{}) error {
	return db.Transaction(func (tx *gorm.DB) error {
		var city domain.City
		if err := tx.First(&city, "board_id = ? and id = ?", boardID, cityID).Error; err != nil {
			return err
		}

		return tx.Delete(&city).Error
	})
}
//
//func AddSpaceToCity(form *CitySpaceForm) (*CitySpace, error) {
//	form.NormalizeInputs()
//
//	if !form.IsValid(tx) {
//		return nil, ErrInvalidForm
//	}
//
//	city, err := GetCityByID(tx, form.CityID)
//	if err != nil {
//		return nil, err
//	}
//
//	space := CitySpace{
//		CityID: city.ID,
//		Order: len(city.CitySpaces),
//		SpaceType: form.SpaceType,
//		RequiredPrivilege: form.RequiredPrivilege,
//	}
//
//	if err = tx.Save(&space).Error; err != nil {
//		return nil, err
//	}
//
//	return &space, nil
//}
