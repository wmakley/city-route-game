package main

import (
	"city-route-game/admin"
	"city-route-game/domain"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

func main() {
	var err error

	var listenAddr string
	var port int
	var migrate bool
	flag.StringVar(&listenAddr, "listenaddr", "", "address to listen on (default \"\")")
	flag.IntVar(&port, "port", 8080, "port to listen on (default 8080)")
	flag.BoolVar(&migrate, "migrate", false, "Migrate database on startup")
	flag.Parse()

	db, err = gorm.Open(sqlite.Open("data/city-route-game.sqlite"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}

	if migrate {
		err = db.AutoMigrate(domain.Models()...)
		if err != nil {
			panic("Error migrating database: " + err.Error())
		}

		fmt.Println("Database migration successful!")
		os.Exit(0)
	}

	admin.Init(db, "./templates")

	router := admin.NewAdminRouter(true)

	listenAddrFull := fmt.Sprintf("%s:%d", listenAddr, port)
	fmt.Println("Listening on", listenAddrFull)
	log.Fatal(http.ListenAndServe(listenAddrFull, router))
}
