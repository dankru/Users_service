package main

import (
	"database/sql"
	"fmt"
	"github.com/dankru/Commissions_simple/internal/repository/psql"
	"github.com/dankru/Commissions_simple/internal/service"
	"github.com/dankru/Commissions_simple/internal/transport/rest"
	hash "github.com/dankru/Commissions_simple/pkg/hasher"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
	fmt.Println("Connecting to DB with DSN:", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("DB init failure: ", err.Error())
	}
	defer db.Close()

	hasher := hash.NewSHA1Hasher("RqijtrEJTQ0wtqTEsGNHrownSaltIGj")

	userRepo := psql.NewRepository(db)
	userService := service.NewService(userRepo, hasher, []byte("Secret here"))
	handler := rest.NewHandler(userService)
	handler.InitRouter()
	if err := http.ListenAndServe(":8080", handler.InitRouter()); err != nil {
		log.Fatal("Failed to run server", err.Error())
	}
}
