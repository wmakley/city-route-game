// +heroku goVersion go1.16
// +heroku install ./cmd/admin/main.go
module city-route-game

go 1.16

require (
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/schema v1.2.0
	gorm.io/driver/postgres v1.1.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.9
)
