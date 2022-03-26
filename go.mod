// +heroku goVersion go1.16
// +heroku install ./cmd/admin/main.go
module city-route-game

go 1.16

require (
	github.com/assertgo/assert v2.0.0+incompatible
	github.com/gorilla/csrf v1.7.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/schema v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/lib/pq v1.3.0
	gorm.io/driver/postgres v1.1.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.9
)
