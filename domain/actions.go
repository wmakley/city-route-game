package domain

import (
	"gorm.io/gorm"
)

func CreateBoard(tx *gorm.DB, form *BoardForm) (err error) {
	form.NormalizeInputs()

	if !form.IsValid(tx) {
		return ErrInvalidForm
	}

	board := Board{
		Name:   form.Name,
		Width:  form.Width,
		Height: form.Height,
	}

	if board.Width == 0 {
		board.Width = 800
	}

	if board.Height == 0 {
		board.Height = 500
	}

	if err := tx.Save(&board).Error; err != nil {
		return err
	}

	return nil
}

func UpdateBoard(tx *gorm.DB, form *BoardForm, board *Board) (err error) {
	form.NormalizeInputs()

	if !form.IsValid(tx) {
		err = ErrInvalidForm
		return
	}

	board.Name = form.Name
	board.Width = form.Width
	board.Height = form.Height

	err = tx.Save(&board).Error
	if err != nil {
		return
	}

	return nil
}

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

func CreateCity(tx *gorm.DB, boardId interface{}, form *CityForm) (city City, err error) {
	form.NormalizeInputs()

	if !form.IsValid() {
		err = ErrInvalidForm
		return
	}

	var board Board
	if err = tx.First(&board, boardId).Error; err != nil {
		return
	}

	city = City{
		BoardID:  board.ID,
		Name:     form.Name,
		Position: form.Position,
	}

	if err = tx.Save(&city).Error; err != nil {
		return
	}

	return city, nil
}

func UpdateCity(tx *gorm.DB, boardId interface{}, form *CityForm) (city City, err error) {
	form.NormalizeInputs()

	if !form.IsValid() {
		err = ErrInvalidForm
		return
	}

	var board Board
	if err = tx.First(&board, boardId).Error; err != nil {
		return
	}

	if err = tx.First(&city, form.ID).Error; err != nil {
		return
	}

	city.Name = form.Name
	city.Position.X = form.Position.X
	city.Position.Y = form.Position.Y

	if err = tx.Save(&city).Error; err != nil {
		return
	}

	return
}

func DeleteCityByBoardIDAndCityID(tx *gorm.DB, boardId interface{}, cityId interface{}) (err error) {
	var city City

	err = tx.First(&city, "board_id = ? AND id = ?", boardId, cityId).Error
	if err != nil {
		return err
	}

	return DeleteCity(tx, &city)
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

func AddSpaceToCity(tx *gorm.DB, form CitySpaceForm, city *City) (err error) {
	form.NormalizeInputs()

	if !form.IsValid(tx) {
		err = ErrInvalidForm
		return
	}

}
