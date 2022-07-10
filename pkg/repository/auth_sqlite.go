package repository

import "database/sql"

type AuthSQLite struct{
	db *sql.DB
}

func NewAuthSQLite(db *sql.DB) *AuthSQLite{
	return &AuthSQLite{db: db}
}

