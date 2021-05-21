package admin

import (
	"city-route-game/domain"
	city2 "city-route-game/domain/city"
	"city-route-game/util"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func ListCitiesHandler(w http.ResponseWriter, r *http.Request) {
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

	var err error
	var cityForm domain.CityForm

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&cityForm); err != nil {
		panic(err)
	}

	var city domain.City

	err = db.Transaction(func(tx *gorm.DB) error {
		city, err = city2.CreateCity(tx, boardId, &cityForm)
		if err != nil {
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
	cityId, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		panic(err)
	}

	var cityForm domain.CityForm

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cityForm); err != nil {
		panic(err)
	}
	cityForm.ID = uint(cityId)

	var city domain.City
	err = db.Transaction(func(tx *gorm.DB) error {
		var err error

		if city, err = city2.UpdateCity(tx, boardId, &cityForm); err != nil {
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

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := city2.DeleteCityByBoardIDAndCityID(tx, boardId, cityId); err != nil {
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
