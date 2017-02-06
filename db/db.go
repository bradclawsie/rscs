// Package db exposes the database interface to rscs. This package can be
// used in conjunction with the rscs daemon or integrated directly into a client
// program.
package db

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	_ "github.com/mattn/go-sqlite3" //
	"io/ioutil"
)

// RscsDB contains the state values for communicating with the underlying sqlite file.
type RscsDB struct {
	sqliteDBFile string
	db           *sql.DB
	dbSHA256     string
	readOnly     bool
}

// dbFileSHA256 calculates the hash of the db file.
func dbFileSHA256(sqliteDBFile string) (string, error) {
	h := sha256.New()
	// The constructor needs to test for the file existing
	// and being readable: the read below suffices.
	contents, fileErr := ioutil.ReadFile(sqliteDBFile)
	if fileErr != nil {
		return "", fileErr
	}
	h.Write(contents)
	return hex.EncodeToString(h.Sum(nil)), nil
}

// NewRscsDB initializes a new RscsDB instance. Write access is set
// with `readOnly` set to true or false; it is assumed that most use will
// be read-only hence the presumption in the parameter name.
func NewRscsDB(sqliteDBFile string, readOnly bool) (*RscsDB, error) {
	dbSHA256, hashErr := dbFileSHA256(sqliteDBFile)
	if hashErr != nil {
		return nil, hashErr
	}
	db, connErr := sql.Open("sqlite3", sqliteDBFile)
	if connErr != nil {
		return nil, connErr
	}
	return &RscsDB{
		sqliteDBFile: sqliteDBFile,
		db:           db,
		dbSHA256:     dbSHA256,
		readOnly:     readOnly}, nil
}

// ReadOnly tells if this DB read-only.
func (r *RscsDB) ReadOnly() bool {
	return r.readOnly
}

// SHA256 returns the hash calculated for the DB file.
func (r *RscsDB) SHA256() string {
	return r.dbSHA256
}
