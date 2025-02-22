package pg_db

import (
	"database/sql"
	"fmt"
	"log"
)

const driver = "postgres"

type Connection struct {
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
}

type PostgresqlDB struct {
	DB *sql.DB
}

func NewPostgreSQLDB(conn Connection) *PostgresqlDB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", conn.DB_HOST, conn.DB_PORT, conn.DB_USER, conn.DB_NAME, conn.DB_PASSWORD)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatal("DB init failure: ", err.Error())
	}

	if err := db.Ping(); err != nil {
		db.Close()
		log.Fatal(err.Error())
	}

	return &PostgresqlDB{DB: db}
}

func (postgres *PostgresqlDB) Close() {
	if err := postgres.DB.Close(); err != nil {
		log.Println("error closing db: ", err.Error())
	}
}
