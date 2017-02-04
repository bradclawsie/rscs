package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
)
	
type RscsDB struct {
	sqliteDBFile string
	db           *sql.DB
	readOnly     bool
}

func NewRscsDB(sqliteDBFile string, readOnly bool) (*RscsDB, error) {
	if _, fileErr := os.Open(sqliteDBFile); fileErr != nil {
		return nil, fileErr
	}
	db, connErr := sql.Open("sqlite3", sqliteDBFile)
	if connErr != nil {
		return nil, connErr
	}
	return &RscsDB{db:db,sqliteDBFile:sqliteDBFile,readOnly:readOnly}, nil
}

