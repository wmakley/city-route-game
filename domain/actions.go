package domain

import "gorm.io/gorm"

func DeleteBoard(tx *gorm.DB, id uint) (err error) {
	var board Board
	if err = tx.First(&board, id).Error; err != nil {
		return
	}

	var cities []City
	if err = tx.Find(&cities, "board_id = ?", board.ID).Error; err != nil {
		return
	}

	for _, city := range cities {
		if err = DeleteCity(tx, &city); err != nil {
			return
		}
	}

	if err = tx.Delete(&board).Error; err != nil {
		return
	}

	return nil
}

func DeleteCity(tx *gorm.DB, city *City) (err error) {
	if err = tx.Delete(&CitySpace{}, "city_id = ?", city.ID).Error; err != nil {
		return
	}

	if err = tx.Delete(&city).Error; err != nil {
		return
	}

	return nil
}
