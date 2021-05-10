package admin

import (
	"city-route-game/domain"
	"city-route-game/util"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func ListCitiesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ListCitiesHandler")

	vars := mux.Vars(r)
	boardId := vars["boardId"]

	var board domain.Board
	var cities []domain.City

	err := db.Transaction(func(tx *gorm.DB) error {
		err := db.First(&board, boardId).Error
		if err != nil {
			return err
		}

		err = db.Where("board_id = ?", boardId).Order("id").Find(&cities).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		handleDBErr(w, r, err)
		return
	}

	util.MustReturnJson(w, cities)
}

func CreateCityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardId := vars["boardId"]

	var cityForm CityForm

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cityForm); err != nil {
		panic(err)
	}

	cityForm.NormalizeInputs()

	var city domain.City

	err := db.Transaction(func(tx *gorm.DB) error {
		var board domain.Board
		if err := db.First(&board, boardId).Error; err != nil {
			return err
		}

		city = domain.City{
			BoardID:  board.ID,
			Name:     cityForm.Name,
			Position: cityForm.Position,
		}

		if err := db.Save(&city).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		handleDBErr(w, r, err)
		return
	}

	util.MustReturnJson(w, &city)
}

func UpdateCityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardId := vars["boardId"]
	cityId := vars["id"]

	var cityForm CityForm

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cityForm); err != nil {
		panic(err)
	}

	cityForm.NormalizeInputs()

	var board domain.Board
	var city domain.City

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := db.First(&board, boardId).Error; err != nil {
			return err
		}

		if err := db.First(&city, cityId).Error; err != nil {
			return err
		}

		city.Name = cityForm.Name
		city.Position.X = cityForm.Position.X
		city.Position.Y = cityForm.Position.Y

		if err := db.Save(&city).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		handleDBErr(w, r, err)
		return
	}

	util.MustReturnJson(w, &city)
}

func DeleteCityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardId := vars["boardId"]
	cityId := vars["id"]

	var city domain.City
	err := db.Transaction(func(tx *gorm.DB) error {
		err := db.First(&city, "board_id = ? AND id = ?", boardId, cityId).Error
		if err != nil {
			return err
		}

		err = db.Delete(&city).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		handleDBErr(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
