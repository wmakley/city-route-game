package admin

import (
	"city-route-game/internal/app"
	"city-route-game/util"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type CityController struct {
	Controller
	boardEditorService app.BoardEditorService
}

func NewCityController(config ControllerConfig, service app.BoardEditorService) CityController {
	return CityController{
		Controller{
			FormDecoder:  config.FormDecoder,
			TemplateRoot: config.TemplateRoot,
			AssetHost:    config.AssetHost,
		},
		service,
	}
}

func (c CityController)Index(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardId := vars["boardId"]

	cities, err := c.boardEditorService.ListCitiesByBoardID(r.Context(), boardId)
	if err != nil {
		c.HandleServiceError(err, w, r)
		return
	}

	util.MustReturnJson(w, cities)
}

func (c CityController)Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardId := vars["boardId"]

	var err error
	var cityForm app.CityForm

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&cityForm); err != nil {
		panic(err)
	}

	city, err := c.boardEditorService.CreateCity(r.Context(), boardId, &cityForm)
	if err != nil {
		c.HandleServiceError(err, w, r)
	}

	util.MustReturnJson(w, &city)
}

func (c CityController)Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//boardId := vars["boardId"]
	cityId := vars["id"]

	var cityForm app.CityForm

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cityForm); err != nil {
		panic(err)
	}

	updatedCity, err := c.boardEditorService.UpdateCity(r.Context(), cityId, &cityForm)
	if err != nil {
		c.HandleServiceError(err, w, r)
		return
	}

	util.MustReturnJson(w, updatedCity)
}

func (c CityController)Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//boardId := vars["boardId"]
	cityId := vars["id"]

	err := c.boardEditorService.DeleteCity(r.Context(), cityId)
	if err != nil {
		c.HandleServiceError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
