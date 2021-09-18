package gorm_board_crud_repository

import (
	"city-route-game/internal/app"
	"context"
	"errors"
	"gorm.io/gorm"
)

func NewGormBoardCrudRepository(db *gorm.DB) app.BoardCrudRepository {
	return &gormBoardRepository{
		db: db,
	}
}

type gormBoardRepository struct {
	db *gorm.DB
}

func (p gormBoardRepository) GetBoardByID(ctx context.Context, id app.ID) (*app.Board, error) {
	var board Board
	if err := p.db.WithContext(ctx).First(&board, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app.NewBoardNotFoundError(id)
		} else {
			return nil, err
		}
	}
	return newDomainBoardFromGormBoard(&board), nil
}

func (p gormBoardRepository) CreateBoard(ctx context.Context, board *app.Board) error {
	var gormBoard *Board
	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var dupe Board
		err := tx.First(&dupe, "name = ?", board.Name).Error
		if err == nil {
			return app.ErrNameTaken
		}
		if !errors.Is(gorm.ErrRecordNotFound, err) {
			return err
		}

		gormBoard, err = newGormBoardFromDomainBoard(board)
		if err != nil {
			return err
		}

		if err := tx.Save(gormBoard).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	board.ID = gormBoard.ID
	board.Name = gormBoard.Name
	board.CreatedAt = gormBoard.CreatedAt
	board.UpdatedAt = gormBoard.UpdatedAt
	board.Width = gormBoard.Width
	board.Height = gormBoard.Height

	return nil
}

func (p gormBoardRepository) UpdateBoard(ctx context.Context, id app.ID, updateFn func(board *app.Board) (*app.Board, error)) (*app.Board, error) {
	var domainBoard *app.Board
	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var board Board
		err := tx.First(&board, id).Error
		if err != nil {
			return err
		}

		domainBoard = newDomainBoardFromGormBoard(&board)
		updatedBoard, err := updateFn(domainBoard)
		if err != nil {
			return err
		}
		if updatedBoard == nil {
			panic("updateFn returned nil board and nil error")
		}

		// check for duplicates
		var dupe Board
		err = tx.First(&dupe, "id <> ? and name = ?", id, updatedBoard.Name).Error
		if err == nil {
			return app.ErrNameTaken
		}
		if !errors.Is(gorm.ErrRecordNotFound, err) {
			return err
		}

		updatedGormBoard, err := newGormBoardFromDomainBoard(updatedBoard)
		if err != nil {
			return err
		}

		err = tx.Save(updatedGormBoard).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return domainBoard, nil
}

func (p gormBoardRepository) ListBoards(ctx context.Context) ([]app.Board, error) {
	var boards []Board
	if err := p.db.WithContext(ctx).Find(&boards).Error; err != nil {
		return nil, err
	}

	domainBoards := make([]app.Board, 0, len(boards))
	for _, board := range boards {
		b := *newDomainBoardFromGormBoard(&board)
		domainBoards = append(domainBoards, b)
	}

	return domainBoards, nil
}

func (p gormBoardRepository) DeleteBoardByID(ctx context.Context, id app.ID) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var board Board
		var err error

		if err = tx.First(&board, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return app.NewBoardNotFoundError(id)
			}
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

func (p gormBoardRepository) ListCitiesByBoardID(ctx context.Context, boardID app.ID) ([]app.City, error) {
	var cities []City
	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var board Board
		err := tx.First(&board, boardID).Error
		if err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				return app.NewBoardNotFoundError(boardID)
			}
			return err
		}

		err = tx.Model(&City{}).Preload("CitySpaces").
			Where("board_id = ?", boardID).
			Find(&cities).
			Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	appCities := make([]app.City, 0, len(cities))
	for _, city := range cities {
		appCities = append(appCities, *newAppCityFromGormCity(&city))
	}

	return appCities, nil
}

func (p gormBoardRepository) GetCityByID(ctx context.Context, id app.ID) (*app.City, error) {
	var city City

	err := p.db.WithContext(ctx).
		Preload("CitySpaces", func(spaces *gorm.DB) *gorm.DB {
			return spaces.Order("`city_spaces`.`order` ASC")
		}).First(&city, id).Error
	if err != nil {
		return nil, err
	}

	return newAppCityFromGormCity(&city), nil
}

func (p gormBoardRepository) CreateCity(ctx context.Context, city *app.City) error {
	gormCity, err := newGormCityFromAppCity(city)
	if err != nil {
		return err
	}

	if err := p.db.WithContext(ctx).Save(gormCity).Error; err != nil {
		return err
	}

	city.ID = gormCity.ID
	city.CreatedAt = gormCity.CreatedAt
	city.UpdatedAt = gormCity.UpdatedAt
	city.Name = gormCity.Name
	city.BoardID = gormCity.BoardID
	city.Position.X = gormCity.Position.X
	city.Position.Y = gormCity.Position.Y

	return nil
}

func (p gormBoardRepository) UpdateCity(ctx context.Context, id app.ID, updateFn func(city *app.City) (*app.City, error)) (*app.City, error) {
	var updatedCity *app.City
	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var city City
		err := tx.First(&city, id).Error
		if err != nil {
			return err
		}

		appCity := newAppCityFromGormCity(&city)

		updatedCity, err = updateFn(appCity)
		if err != nil {
			return err
		}
		if updatedCity == nil {
			panic("updateFn returned nil error and nil city")
		}

		updatedGormCity, err := newGormCityFromAppCity(updatedCity)
		if err != nil {
			return err
		}

		err = tx.Save(updatedGormCity).Error
		if err != nil {
			return err
		}

		updatedCity.UpdatedAt = updatedGormCity.UpdatedAt

		return nil
	})
	if err != nil {
		return nil, err
	}

	return updatedCity, nil
}

func (p gormBoardRepository) DeleteCityByBoardIDAndCityID(ctx context.Context, boardID app.ID, cityID app.ID) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var city app.City
		if err := tx.First(&city, "board_id = ? AND id = ?", boardID, cityID).Error; err != nil {
			return err
		}

		if err := tx.Delete(&app.CitySpace{}, "city_id = ?", cityID).Error; err != nil {
			return err
		}

		if err := tx.Delete(city).Error; err != nil {
			return err
		}

		return nil
	})
}

func (p gormBoardRepository) DeleteCityByID(ctx context.Context, id app.ID) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var city City
		err := tx.First(&city, id).Error
		if err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				return app.NewRecordNotFoundError("City", id)
			}
			return err
		}

		if err = tx.Delete(&city).Error; err != nil {
			return err
		}

		return nil
	})
}

func (p gormBoardRepository) CreateCitySpace(ctx context.Context, citySpace *app.CitySpace) error {
	gormSpace, err := newGormCitySpaceFromAppCitySpace(citySpace)
	if err != nil {
		return err
	}

	err = p.db.WithContext(ctx).Save(gormSpace).Error
	if err != nil {
		return err
	}

	citySpace.ID = gormSpace.ID
	citySpace.CreatedAt = gormSpace.CreatedAt
	citySpace.UpdatedAt = gormSpace.UpdatedAt
	citySpace.Order = gormSpace.Order
	citySpace.RequiredPrivilege = gormSpace.RequiredPrivilege
	citySpace.SpaceType = gormSpace.SpaceType

	return nil
}

func (p gormBoardRepository) UpdateCitySpace(ctx context.Context, id app.ID, updateFn func(*app.CitySpace) (*app.CitySpace, error)) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var citySpace *app.CitySpace

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

func (p gormBoardRepository) GetCitySpacesByCityID(ctx context.Context, cityID app.ID) ([]app.CitySpace, error) {
	var spaces []CitySpace
	if err := p.db.WithContext(ctx).
		Find(&spaces, "city_id = ?", cityID).
		Error; err != nil {
		return nil, err
	}

	appSpaces := make([]app.CitySpace, 0, len(spaces))
	for _, space := range spaces {
		s := *newAppCitySpaceFromGormCitySpace(&space)
		appSpaces = append(appSpaces, s)
	}

	return appSpaces, nil
}

func (p gormBoardRepository) DeleteCitySpaceByID(ctx context.Context, id app.ID) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var space CitySpace
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
