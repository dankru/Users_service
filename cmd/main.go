package main

import (
	"github.com/dankru/Commissions_simple/internal/grpc"
	"github.com/dankru/Commissions_simple/internal/repository/pg_repo"
	"github.com/dankru/Commissions_simple/internal/server"
	"github.com/dankru/Commissions_simple/internal/service"
	"github.com/dankru/Commissions_simple/internal/transport/rest"
	"github.com/dankru/Commissions_simple/pkg/database/pg_db"
	hash "github.com/dankru/Commissions_simple/pkg/hasher"
	_ "github.com/lib/pq"
	"os"
	"time"
)

func main() {
	grpc.ParseToken()

	conn := pg_db.Connection{
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_NAME:     os.Getenv("DB_NAME"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
	}

	postgres := pg_db.NewPostgreSQLDB(conn)
	defer postgres.Close()

	hasher := hash.NewSHA1Hasher(os.Getenv("SALT"))

	userRepo := pg_repo.NewRepository(postgres.DB)
	authRepo := pg_repo.NewAuthRepository(postgres.DB)
	tokensRepo := pg_repo.NewTokens(postgres.DB)

	userService := service.NewService(userRepo)
	authService := service.NewAuthService(authRepo, tokensRepo, hasher, []byte(os.Getenv("HMAC_SECRET")))

	handler := rest.NewHandler(authService, userService)

	srv := server.NewServer("0.0.0.0:8080", time.Second*15, time.Second*15, time.Second*60, handler.InitRouter())

	srv.Run()

}
