package admin

import (
	"city-route-game/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAdminRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(
		middleware.RequestLogger,
		middleware.RecoverPanic,
		middleware.CSRFMitigation,
		middleware.ParseFormData,
		middleware.HtmlContentType,
		middleware.PreventCache,
	)

	router.Handle("/", http.RedirectHandler("/boards", http.StatusFound))

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

	router.Handle("/{file}", http.FileServer(http.Dir("static/admin")))

	return router
}
