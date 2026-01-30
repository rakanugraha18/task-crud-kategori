package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // SQLite driver
)

func InitDB(dbPath string) (*sql.DB, error) {
	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Penting: cek koneksi
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Set connection pool settings (optional tapi recommended)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return db, nil
}
