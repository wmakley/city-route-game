package main

import (
	"flag"
	"fmt"
	"hansa/admin"
	"hansa/domain"
	"hansa/gamesapi"
	"hansa/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func main() {
	var err error

	var listenAddr string
	var port int
	flag.StringVar(&listenAddr, "listenaddr", "", "address to listen on (default \"\")")
	flag.IntVar(&port, "port", 8080, "port to listen on (default 8080)")
	flag.Parse()

	db, err = gorm.Open(sqlite.Open("./hansa.sqlite"), &gorm.Config{})
	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}

	err = db.AutoMigrate(&domain.Game{}, &domain.Board{}, &domain.Player{}, &domain.PlayerBoard{}, &domain.PlayerBonusToken{}, &domain.BonusToken{}, &domain.RouteBonusToken{}, &domain.City{}, &domain.CitySlot{}, &domain.Route{}, &domain.RouteSlot{})
	if err != nil {
		panic("Error migrating database: " + err.Error())
	}

	admin.Init(db)
	gamesapi.Init(db)

	router := mux.NewRouter().StrictSlash(true)
	router.Use(
		middleware.RequestLogger,
		middleware.RecoverPanic,
	)

	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.CorsHeaders, middleware.JsonContentType)
	api.HandleFunc("/games", gamesapi.ListGamesHandler).Methods("GET", "OPTIONS")
	api.HandleFunc("/games", gamesapi.CreateGameHandler).Methods("POST", "OPTIONS")
	api.HandleFunc("/games/{id}", gamesapi.GetGameByIDHandler).Methods("GET", "OPTIONS")

	adminR := router.PathPrefix("/admin").Subrouter()
	adminR.Use(
		middleware.CSRFMitigation,
		middleware.ParseFormData,
		middleware.HtmlContentType,
		middleware.PreventCache,
	)
	adminR.Handle("/", http.RedirectHandler("/admin/boards", http.StatusFound))
	adminR.HandleFunc("/boards", admin.BoardsIndexHandler).Methods("GET")
	adminR.HandleFunc("/boards/new", admin.NewBoardHandler).Methods("GET")
	adminR.HandleFunc("/boards", admin.CreateBoardHandler).Methods("POST")
	adminR.HandleFunc("/boards/{id}", admin.GetBoardByIdHandler).Methods("GET")
	adminR.HandleFunc("/boards/{id}/edit", admin.EditBoardHandler).Methods("GET")
	adminR.HandleFunc("/boards/{id}", admin.UpdateBoardHandler).Methods("POST", "PATCH", "PUT")
	adminR.HandleFunc("/boards/{id}", admin.DeleteBoardHandler).Methods("DELETE")

	router.Handle("/{file}", http.FileServer(http.Dir("static")))

	listenAddrFull := fmt.Sprintf("%s:%d", listenAddr, port)
	fmt.Println("Listening on", listenAddrFull)
	log.Fatal(http.ListenAndServe(listenAddrFull, router))
}
