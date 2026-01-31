package database

import "database/sql"

// =======================
// MAIN MIGRATION
// =======================
func Migrate(db *sql.DB) error {
	if err := migrateCategories(db); err != nil {
		return err
	}

	if err := migrateProducts(db); err != nil {
		return err
	}

	return nil
}

// =======================
// MIGRATE CATEGORIES
// =======================
func migrateCategories(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT
	);
	`
	_, err := db.Exec(query)
	return err
}

// =======================
// MIGRATE PRODUCTS
// =======================
func migrateProducts(db *sql.DB) error {
	// 1. ensuere products table exists
	createProducts := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price INTEGER NOT NULL,
		stock INTEGER NOT NULL
	);
	`
	if _, err := db.Exec(createProducts); err != nil {
		return err
	}

	// 2. check if category_id column exists
	hasCategoryID, err := columnExists(db, "products", "category_id")
	if err != nil {
		return err
	}

	// 3. add category_id column if not exists
	if !hasCategoryID {
		alter := `
		ALTER TABLE products
		ADD COLUMN category_id INTEGER;
		`
		if _, err := db.Exec(alter); err != nil {
			return err
		}
	}

	return nil
}

// =======================
// CHECK COLUMN EXISTS
// =======================

func columnExists(db *sql.DB, tableName, columnName string) (bool, error) {
	rows, err := db.Query(`PRAGMA table_info(` + tableName + `);`)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name string
		var ctype string
		var notnull int
		var dflt sql.NullString
		var pk int

		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return false, err
		}

		if name == columnName {
			return true, nil
		}
	}

	return false, nil
}
