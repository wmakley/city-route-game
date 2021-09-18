package main

import (
	"city-route-game/admin"
	"city-route-game/internal/app"
	"city-route-game/internal/gorm_board_crud_repository"
	"flag"
	"fmt"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	flag.BoolVar(&migrate, "migrate", false, "Migrate gorm_board_crud_repository on startup")
	flag.StringVar(&assetHost, "assethost", "", "Optional asset host domain")
	flag.StringVar(&databaseUrl, "gorm_board_crud_repository-url", "host=localhost user=william password=password dbname=hansa_dev port=5432 sslmode=disable TimeZone=UTC", "Database URL")
	flag.StringVar(&ipWhitelist, "ipwhitelist", "", "Optional IP Whitelist")
	flag.Parse()

	fmt.Println("Database URL:", databaseUrl)
	var db *gorm.DB
	db, err = gorm.Open(postgres.Open(databaseUrl), &gorm.Config{
		DisableNestedTransaction: true,
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic("Error connecting to gorm_board_crud_repository: " + err.Error())
	}

	if migrate {
		err = db.AutoMigrate(gorm_board_crud_repository.Models()...)
		if err != nil {
			panic("Error migrating gorm_board_crud_repository: " + err.Error())
		}

		fmt.Println("Database migration successful!")
	}

	var splitIPs []string
	if ipWhitelist != "" {
		splitIPs = strings.Split(ipWhitelist, ",")
		fmt.Printf("Whitelisted IPs: %+v\n", splitIPs)
	}

	boardRepo := gorm_board_crud_repository.NewGormBoardCrudRepository(db)
	boardEditorService := app.NewBoardEditorService(boardRepo)

	controllerConfig := admin.ControllerConfig{
		FormDecoder: schema.NewDecoder(),
		TemplateRoot: "./templates",
		AssetHost: "",
	}

	boardController := admin.NewBoardController(controllerConfig, boardEditorService)
	cityController := admin.NewCityController(controllerConfig, boardEditorService)

	router := admin.NewAdminRouter(&boardController, &cityController, splitIPs, true)

	listenAddrFull := fmt.Sprintf("%s:%d", listenAddr, port)
	fmt.Println("Listening on", listenAddrFull)
	log.Fatal(http.ListenAndServe(listenAddrFull, router))
}
