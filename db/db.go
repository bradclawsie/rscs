// Package db exposes the database interface to rscs. This package can be
// used in conjunction with the rscs daemon or integrated directly into a client
// program.
package db

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3" //
	"io/ioutil"
)

const (
	// KVTableName is the KV table name.
	KVTableName = "kv"
	// KVPrimaryKeyColumn is the KV primary key.
	KVPrimaryKeyColumn = "key"
	// KVValueColumn is the KV value.
	KVValueColumn = "value"
)

// RscsDB contains the state values for communicating with the underlying sqlite file.
type RscsDB struct {
	sqliteDBFile string
	db           *sql.DB
	dbSHA256     string
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

// NewRscsDB initializes a new RscsDB instance.
func NewRscsDB(sqliteDBFile string) (*RscsDB, error) {
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
		dbSHA256:     dbSHA256}, nil
}

// SHA256 returns the hash calculated for the DB file.
func (r *RscsDB) SHA256() string {
	return r.dbSHA256
}

// CreateTable will create a new kv table. You will need to DROP independently if needed.
func (r *RscsDB) CreateTable() error {
	queryStr := fmt.Sprintf("CREATE TABLE %s (%s VARCHAR(255) PRIMARY KEY, %s TEXT NOT NULL)",
		KVTableName, KVPrimaryKeyColumn, KVValueColumn)
	_, createErr := r.db.Exec(queryStr)
	return createErr
}

// DropTable will drop the kv table.
func (r *RscsDB) DropTable() error {
	queryStr := fmt.Sprintf("DROP TABLE %s", KVTableName)
	_, dropErr := r.db.Exec(queryStr)
	return dropErr
}

// Insert will insert a new key/value pair.
func (r *RscsDB) Insert(key, value string) (int, error) {
	if key == "" || value == "" {
		return 0, errors.New("insert empty key or value")
	}
	queryStr := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES ($1, $2)",
		KVTableName, KVPrimaryKeyColumn, KVValueColumn)
	result, insertErr := r.db.Exec(queryStr, key, value)
	if insertErr != nil {
		return 0, insertErr
	}
	rowCount, rowCountErr := result.RowsAffected()
	if rowCountErr != nil {
		return 0, rowCountErr
	}
	return int(rowCount), nil
}

// Delete will delete a row with ID key.
func (r *RscsDB) Delete(key string) (int, error) {
	if key == "" {
		return 0, errors.New("delete empty key")
	}
	queryStr := fmt.Sprintf("DELETE FROM %s where %s = $1",
		KVTableName, KVPrimaryKeyColumn)
	result, deleteErr := r.db.Exec(queryStr, key)
	if deleteErr != nil {
		return 0, deleteErr
	}
	rowCount, rowCountErr := result.RowsAffected()
	if rowCountErr != nil {
		return 0, rowCountErr
	}
	return int(rowCount), nil
}

// Get returns the value string for the key string. The second return
// value is a 'found' flag that easily distinguishes a db error case
// from that of no matching row.
func (r *RscsDB) Get(key string) (string, bool, error) {
	if key == "" {
		return "", false, errors.New("key is an empty string")
	}
	queryStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s=?",
		KVValueColumn, KVTableName, KVPrimaryKeyColumn)
	var value string
	selectErr := r.db.QueryRow(queryStr, key).Scan(&value)
	switch {
	case selectErr == sql.ErrNoRows:
		return "", false, nil
	case selectErr != nil:
		return "", false, selectErr
	default:
		return value, true, nil
	}
}
