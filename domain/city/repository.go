package city

import (
	"city-route-game/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Transaction(func (repository Repository) error) error
	ListCitiesByBoardID(boardID interface{}) ([]domain.City, error)
	GetCityByID(id interface{}) (*domain.City, error)
	SaveCity(city *domain.City) error
	DeleteCityByBoardIDAndCityID(boardID interface{}, cityID interface{}) error
	DeleteCity(city *domain.City) error
}

type GormRepositoryImpl struct {
	db *gorm.DB
}

func NewGormRepositoryImpl(db *gorm.DB) *GormRepositoryImpl {
	return &GormRepositoryImpl{
		db: db,
	}
}

func (p *GormRepositoryImpl)Transaction(callback func (Repository) error) error {
	return p.db.Transaction(func (tx *gorm.DB) error {
		return callback(NewGormRepositoryImpl(tx))
	})
}

func (p *GormRepositoryImpl)ListCitiesByBoardID(boardID interface{}) ([]domain.City, error) {
	var cities []domain.City
	err := p.db.Where("board_id = ?", boardID).Preload("CitySpaces").Find(&cities).Error
	if err != nil {
		return nil, err
	}
	return cities, nil
}

func (p *GormRepositoryImpl)GetCityByID(id interface{}) (*domain.City, error) {
	var city domain.City

	err := p.db.Preload("CitySpaces").Order("city_spaces.order").First(&city, id).Error
	if err != nil {
		return nil, err
	}

	return &city, nil
}

func (p *GormRepositoryImpl)SaveCity(city *domain.City) error {
	if err := p.db.Save(city).Error; err != nil {
		return err
	}
	return nil
}

func (p *GormRepositoryImpl)DeleteCityByBoardIDAndCityID(boardID interface{}, cityID interface{}) error {
	return p.db.Transaction(func (tx *gorm.DB) error {
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

func (p *GormRepositoryImpl)DeleteCity(city *domain.City) error {
	return p.db.Transaction(func (tx *gorm.DB) error {

		if err := tx.Delete(city).Error; err != nil {
			return err
		}

		return nil
	})
}
