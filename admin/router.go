package admin

import (
	"city-route-game/middleware"
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
		ip := getIP(r)
		log.Println("Request IP", ip)

		_, ipFound := config.IPWhitelist[ip]

		if ipFound {
			next.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
