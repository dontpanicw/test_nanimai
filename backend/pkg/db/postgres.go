package db

import "database/sql"

func CreateBalanceTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS balance (
		user_id UUID PRIMARY KEY,
		current INTEGER,
		maximum FLOAT
	);`
	_, err := db.Exec(query)
	return err
}

func CreateUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		user_id UUID PRIMARY KEY,
		login TEXT,
		password TEXT
	);`
	_, err := db.Exec(query)
	return err
}

func CreateServicesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS services (
		service_id    UUID PRIMARY KEY,
		name          TEXT UNIQUE NOT NULL,
		api_key       TEXT NOT NULL,
		created_at    TIMESTAMPTZ DEFAULT now()
	);`
	_, err := db.Exec(query)
	return err
}
