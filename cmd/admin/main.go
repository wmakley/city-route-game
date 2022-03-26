package main

import (
	"city-route-game/internal/admin"
	"city-route-game/internal/app"
	"city-route-game/internal/gorm_board_crud_repository"
	"flag"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/schema"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}

	var listenAddr string
	var port int
	var migrate bool
	var assetHost string
	var ipWhitelist string
	flag.StringVar(&listenAddr, "listenaddr", "", "address to listen on (default \"\")")
	flag.IntVar(&port, "port", 8080, "port to listen on (default 8080)")
	flag.BoolVar(&migrate, "migrate", false, "Migrate gorm_board_crud_repository on startup")
	flag.StringVar(&assetHost, "assethost", "", "Optional asset host domain")
	flag.StringVar(&ipWhitelist, "ipwhitelist", "", "Optional IP Whitelist")
	flag.Parse()

	databaseUrl := os.Getenv("DATABASE_URL")

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
	CSRF := csrf.Protect([]byte("32-byte-long-auth-key"))

	listenAddrFull := fmt.Sprintf("%s:%d", listenAddr, port)
	fmt.Println("Listening on", listenAddrFull)
	log.Fatal(http.ListenAndServe(listenAddrFull, CSRF(router)))
}
