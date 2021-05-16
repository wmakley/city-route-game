package main

import (
	"city-route-game/admin"
	"city-route-game/domain"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

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
	var ipWhitelist string
	flag.StringVar(&listenAddr, "listenaddr", "", "address to listen on (default \"\")")
	flag.IntVar(&port, "port", 8080, "port to listen on (default 8080)")
	flag.BoolVar(&migrate, "migrate", false, "Migrate database on startup")
	flag.StringVar(&assetHost, "assethost", "", "Optional asset host domain")
	flag.StringVar(&databaseUrl, "database-url", "host=localhost user=william password=password dbname=hansa_dev port=5432 sslmode=disable TimeZone=UTC", "Database URL")
	flag.StringVar(&ipWhitelist, "ipwhitelist", "", "Optional IP Whitelist")
	flag.Parse()

	fmt.Println("Database URL:", databaseUrl)
	db, err = gorm.Open(postgres.Open(databaseUrl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
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
	}

	var splitIPs []string
	if ipWhitelist != "" {
		splitIPs = strings.Split(ipWhitelist, ",")
		fmt.Printf("Whitelisted IPs: %+v\n", splitIPs)
	}

	admin.Init(db, "./templates", assetHost, splitIPs)

	router := admin.NewAdminRouter(true)

	listenAddrFull := fmt.Sprintf("%s:%d", listenAddr, port)
	fmt.Println("Listening on", listenAddrFull)
	log.Fatal(http.ListenAndServe(listenAddrFull, router))
}
