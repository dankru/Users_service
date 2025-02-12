package main

import (
	"github.com/dankru/Commissions_simple/internal/repository/psql"
	"github.com/dankru/Commissions_simple/internal/server"
	"github.com/dankru/Commissions_simple/internal/service"
	"github.com/dankru/Commissions_simple/internal/transport/rest"
	"github.com/dankru/Commissions_simple/pkg/database/pgsql"
	hash "github.com/dankru/Commissions_simple/pkg/hasher"
	_ "github.com/lib/pq"
	"os"
	"time"
)

func main() {
	conn := pgsql.Connection{
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_NAME:     os.Getenv("DB_NAME"),
	}
	postgres := pgsql.NewPostgreSQLDB(conn)
	defer postgres.Close()

	hasher := hash.NewSHA1Hasher(os.Getenv("SALT"))

	userRepo := psql.NewRepository(postgres.DB)
	authRepo := psql.NewAuthRepository(postgres.DB)
	tokensRepo := psql.NewTokens(postgres.DB)

	userService := service.NewService(userRepo)
	authService := service.NewAuthService(authRepo, tokensRepo, hasher, []byte(os.Getenv("HMAC_SECRET")))

	handler := rest.NewHandler(authService, userService)

	srv := server.NewServer("0.0.0.0:8080", time.Second*15, time.Second*15, time.Second*60, handler.InitRouter())

	srv.Run()

}
