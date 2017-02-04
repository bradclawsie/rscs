// Package db exposes the database interface to rscs. This package can be
// used in conjunction with the rscs daemon or integrated directly into a client
// program.
package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" //
	"os"
)

// RscsDB contains the state values for communicating with the underlying sqlite file.
type RscsDB struct {
	sqliteDBFile string
	db           *sql.DB
	readOnly     bool
}

// NewRscsDB initializes a new RscsDB instance.
func NewRscsDB(sqliteDBFile string, readOnly bool) (*RscsDB, error) {
	if _, fileErr := os.Open(sqliteDBFile); fileErr != nil {
		return nil, fileErr
	}
	db, connErr := sql.Open("sqlite3", sqliteDBFile)
	if connErr != nil {
		return nil, connErr
	}
	return &RscsDB{db: db, sqliteDBFile: sqliteDBFile, readOnly: readOnly}, nil
}
