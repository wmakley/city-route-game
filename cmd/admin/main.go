package main

import (
	"city-route-game/admin"
	"city-route-game/domain"
	"city-route-game/middleware"
	"flag"
	"fmt"
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

	db, err = gorm.Open(sqlite.Open("data/city-route-game.sqlite"), &gorm.Config{})
	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}

	err = db.AutoMigrate(&domain.Game{}, &domain.Board{}, &domain.Player{}, &domain.PlayerBoard{}, &domain.PlayerBonusToken{}, &domain.BonusToken{}, &domain.RouteBonusToken{}, &domain.City{}, &domain.CitySlot{}, &domain.Route{}, &domain.RouteSlot{})
	if err != nil {
		panic("Error migrating database: " + err.Error())
	}

	admin.Init(db)

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
	router.HandleFunc("/boards", admin.BoardsIndexHandler).Methods("GET")
	router.HandleFunc("/boards/new", admin.NewBoardHandler).Methods("GET")
	router.HandleFunc("/boards", admin.CreateBoardHandler).Methods("POST")
	router.HandleFunc("/boards/{id}", admin.GetBoardByIdHandler).Methods("GET")
	router.HandleFunc("/boards/{id}/edit", admin.EditBoardHandler).Methods("GET")
	router.HandleFunc("/boards/{id}", admin.UpdateBoardHandler).Methods("POST", "PATCH", "PUT")
	router.HandleFunc("/boards/{id}", admin.DeleteBoardHandler).Methods("DELETE")

	router.Handle("/{file}", http.FileServer(http.Dir("static/admin")))

	listenAddrFull := fmt.Sprintf("%s:%d", listenAddr, port)
	fmt.Println("Listening on", listenAddrFull)
	log.Fatal(http.ListenAndServe(listenAddrFull, router))
}
