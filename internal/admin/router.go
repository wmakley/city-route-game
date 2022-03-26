package admin

import (
	"city-route-game/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAdminRouter(
	boardController *BoardController,
	cityController *CityController,
	ipWhitelist []string,
	logRequests bool,
) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	if logRequests {
		router.Use(middleware.RequestLogger)
	}
	if len(ipWhitelist) > 0 {
		router.Use(middleware.NewIPWhiteListMiddleware(ipWhitelist, true))
	}
	router.Use(
		middleware.RecoverPanic,
		middleware.CSRFMitigation,
		middleware.ParseFormData,
		middleware.HtmlContentType,
		middleware.PreventCache,
	)

	router.Handle("/", http.RedirectHandler("/boards/", http.StatusFound))

	boards := router.PathPrefix("/boards").Subrouter()
	boards.HandleFunc("/", boardController.Index).Methods("GET")
	boards.HandleFunc("/new", boardController.New).Methods("GET")
	boards.HandleFunc("/", boardController.Create).Methods("POST")
	boards.HandleFunc("/{id}", boardController.GetById).Methods("GET")
	boards.HandleFunc("/{id}/edit", boardController.Edit).Methods("GET")
	boards.HandleFunc("/{id}", boardController.Update).Methods("POST", "PATCH")
	boards.HandleFunc("/{id}", boardController.Delete).Methods("DELETE")

	cities := boards.PathPrefix("/{boardId}/cities").Subrouter()
	cities.HandleFunc("/", cityController.Index).Methods("GET")
	cities.HandleFunc("/", cityController.Create).Methods("POST")
	cities.HandleFunc("/{id}", cityController.Update).Methods("PUT")
	cities.HandleFunc("/{id}", cityController.Delete).Methods("DELETE")

	router.Handle("/{file}", http.FileServer(http.Dir("static/admin")))

	return router
}
