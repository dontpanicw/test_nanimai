package postgres

import "database/sql"

type PostgresStorageUser struct {
	db *sql.DB
}

func (su *PostgresStorageUser) GetDb() *sql.DB {
	return su.db
}

func NewUserPostgresStorageUser(connStr string) (*PostgresStorageUser, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			db.Close()
		}
	}()
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &PostgresStorageUser{db: db}, nil
}
