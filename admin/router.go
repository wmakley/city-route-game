package admin

import (
	"city-route-game/middleware"
	"city-route-game/util"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAdminRouter(logRequests bool) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	if logRequests {
		router.Use(middleware.RequestLogger)
	}
	if len(config.IPWhitelist) > 0 {
		router.Use(IPWhitelistMiddleware)
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

func IPWhitelistMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, err := util.GetIP(r)
		if err != nil {
			log.Printf("Error getting request IP: %+v\n", err)
			ip = ""
		} else {
			log.Println("Request IP", ip)
		}

		_, ipFound := config.IPWhitelist[ip]

		if ipFound {
			next.ServeHTTP(w, r)
		} else {
			log.Println("IP not found in whitelist")
			http.NotFound(w, r)
		}
	})
}
