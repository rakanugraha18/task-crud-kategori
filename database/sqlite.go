package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // using modernc.org/sqlite driver for SQLite becouse it's pure Go implementation and cross-platform no need CGO
)

func InitDB(dbPath string) (*sql.DB, error) {
	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	// Handle error
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	// Log successful connection
	log.Println("Database connected successfully")
	return db, nil
}
