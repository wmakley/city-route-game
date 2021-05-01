package gamesapi

import (
	"hansa/domain"
	"hansa/util"
	"net/http"

	"github.com/gorilla/mux"
)

func ListGamesHandler(w http.ResponseWriter, r *http.Request) {
	var games []domain.Game

	if err := db.Find(&games).Error; err != nil {
		util.JsonInternalServerError(err.Error(), w, r)
		return
	}

	util.MustEncode(w, games)
}

func GetGameByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	var game domain.Game

	if err := db.Find(&game, key).Error; err != nil {
		util.JsonHandleDbErr(err, w, r)
		return
	}

	util.MustEncode(w, game)
}

func CreateGameHandler(w http.ResponseWriter, r *http.Request) {

}
