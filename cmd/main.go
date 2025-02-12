package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/dankru/Commissions_simple/internal/repository/psql"
	"github.com/dankru/Commissions_simple/internal/service"
	"github.com/dankru/Commissions_simple/internal/transport/rest"
	hash "github.com/dankru/Commissions_simple/pkg/hasher"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

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
	authRepo := psql.NewAuthRepository(db)
	tokensRepo := psql.NewTokens(db)

	userService := service.NewService(userRepo)
	authService := service.NewAuthService(authRepo, tokensRepo, hasher, []byte("Secret here"))

	handler := rest.NewHandler(authService, userService)

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler.InitRouter(), // Pass our instance of gorilla/mux in.
	}
	go func() {
		if err := http.ListenAndServe(":8080", handler.InitRouter()); err != nil {
			log.Fatal("Failed to run server", err.Error())
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
