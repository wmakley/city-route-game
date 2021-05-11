package domain

import "gorm.io/gorm"

func DeleteBoard(tx *gorm.DB, id uint) error {
	var board Board
	if err := tx.First(&board, id).Error; err != nil {
		return err
	}

	var cities []City
	if err := tx.Find(&cities, "board_id = ?", board.ID).Error; err != nil {
		return err
	}

	if err := tx.Delete(&cities).Error; err != nil {
		return err
	}

	if err := tx.Delete(&board).Error; err != nil {
		return err
	}

	return nil
}
