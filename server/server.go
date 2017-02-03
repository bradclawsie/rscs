package server

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type Rscs struct {
	sqliteDBFile string
	db           *sql.DB
}

func NewRscs(sqliteDBFile string) (*Rscs, error) {
	if _, fileErr := os.Open(sqliteDBFile); fileErr != nil {
		return nil, fileErr
	}
	db, connErr := sql.Open("sqlite3", sqliteDBFile)
	if connErr != nil {
		return nil, connErr
	}
	
	// Validate the db

	return &Rscs{db:db,sqliteDBFile:sqliteDBFile}, nil
}
