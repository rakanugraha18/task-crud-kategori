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

	if err := migrationTransaction(db); err != nil {
		return err
	}
	if err := migrationReport(db); err != nil {
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
	// 1. ensure products table exists
	createProducts := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price INTEGER NOT NULL,
		stock INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
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

//======================
// MIGRATION TRANSACTION
//======================
func migrationTransaction(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY
	)`)
	if err != nil {
		return err
	}

	var count int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM schema_migrations WHERE version = '001'
	`).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // sudah di-migrate
	}

	migrationSQL := `
	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		total_amount INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS transaction_details (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		transaction_id INTEGER NOT NULL,
		product_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		subtotal INTEGER NOT NULL,
		FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products(id)
	);
	`

	_, err = db.Exec(migrationSQL)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO schema_migrations (version) VALUES ('001')
	`)
	return err
}

// =======================
// MIGRATE REPORTS
// =======================
func migrationReport(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY
	)`)
	if err != nil {
		return err
	}

	var count int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM schema_migrations WHERE version = '001_report'
	`).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	migrationSQL := `
	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		total_amount INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS transaction_details (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		transaction_id INTEGER NOT NULL,
		product_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		subtotal INTEGER NOT NULL,
		FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products(id)
	);
	`

	if _, err := db.Exec(migrationSQL); err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO schema_migrations (version) VALUES ('001_report')
	`)
	return err
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
