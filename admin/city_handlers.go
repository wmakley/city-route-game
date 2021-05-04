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
	log.Println("CreateCityHandler")

	vars := mux.Vars(r)
	boardId := vars["boardId"]

	var cityForm CityForm

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cityForm); err != nil {
		panic(err)
	}

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
