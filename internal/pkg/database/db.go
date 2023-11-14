package database

import (
	"database/sql"
	"log"
	"os"

	// We are using the pgx driver to connect to PostgreSQL
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Env holds database connection to Postgres
type Env struct {
	DB *sql.DB
}

// ConnectDB tries to connect DB and on succcesful it returns
// DB connection string and nil error, otherwise return empty DB and the corresponding error.
func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", os.Getenv("RDS_DSN"))
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return &sql.DB{}, err
	}
	return db, nil
}
