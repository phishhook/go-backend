package config

import (
	"database/sql"
	"log"

	// We are using the pgx driver to connect to PostgreSQL
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Env holds database connection to Postgres
type Env struct {
	DB *sql.DB
}

// ConnectDB tries to connect DB and on succcesful it returns
// DB connection string and nil error, otherwise return empty DB and the corresponding error.
func ConnectDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Printf("failed to ping the database: %v", err)
		return nil, err
	}
	return db, nil
}
