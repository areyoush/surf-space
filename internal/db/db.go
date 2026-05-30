package db

import (
	"log"
	"database/sql"
	_ "github.com/lib/pq"
)

func Connect(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open database connection: ", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("failed to ping database, check credentials and host: ", err)
	}

	return db
}