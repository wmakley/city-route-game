package admin

import (
	"city-route-game/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAdminRouter(logRequests bool) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	if logRequests {
		router.Use(middleware.RequestLogger)
	}
	if len(config.IPWhitelist) > 0 {
		router.Use(middleware.NewIPWhiteListMiddleware(config.IPWhitelist, true))
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
	boards.HandleFunc("/", BoardsIndexHandler).Methods("GET")
	boards.HandleFunc("/new", NewBoardHandler).Methods("GET")
	boards.HandleFunc("/", CreateBoardHandler).Methods("POST")
	boards.HandleFunc("/{id}", GetBoardByIdHandler).Methods("GET")
	boards.HandleFunc("/{id}/edit", EditBoardHandler).Methods("GET")
	boards.HandleFunc("/{id}", UpdateBoardHandler).Methods("POST", "PATCH")
	boards.HandleFunc("/{id}", DeleteBoardHandler).Methods("DELETE")

	cities := boards.PathPrefix("/{boardId}/cities").Subrouter()
	cities.HandleFunc("/", ListCitiesHandler).Methods("GET")
	cities.HandleFunc("/", CreateCityHandler).Methods("POST")
	cities.HandleFunc("/{id}", UpdateCityHandler).Methods("PUT")
	cities.HandleFunc("/{id}", DeleteCityHandler).Methods("DELETE")

	router.Handle("/{file}", http.FileServer(http.Dir("static/admin")))

	return router
}
