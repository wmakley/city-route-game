package main

import (
	"city-route-game/admin"
	"city-route-game/domain"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
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
	var assetHost string
	var databaseUrl string
	flag.StringVar(&listenAddr, "listenaddr", "", "address to listen on (default \"\")")
	flag.IntVar(&port, "port", 8080, "port to listen on (default 8080)")
	flag.BoolVar(&migrate, "migrate", false, "Migrate database on startup")
	flag.StringVar(&assetHost, "assethost", "", "Optional asset host domain")
	flag.StringVar(&databaseUrl, "database-url", "host=localhost user=william dbname=hansa_dev port=5432 sslmode=disable TimeZone=UTC", "Database URL")
	flag.Parse()

	log.Println("Database URL:", databaseUrl)
	db, err = gorm.Open(postgres.Open(databaseUrl), &gorm.Config{
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

	admin.Init(db, "./templates", assetHost)

	router := admin.NewAdminRouter(true)

	listenAddrFull := fmt.Sprintf("%s:%d", listenAddr, port)
	fmt.Println("Listening on", listenAddrFull)
	log.Fatal(http.ListenAndServe(listenAddrFull, router))
}
