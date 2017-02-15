// Package db exposes the database interface to rscs. This package can be
// used in conjunction with the rscs daemon or integrated directly into a client
// program.
package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3" //
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
}

// NewRscsDB initializes a new RscsDB instance.
func NewRscsDB(sqliteDBFile string) (*RscsDB, error) {
	db, connErr := sql.Open("sqlite3", sqliteDBFile)
	if connErr != nil {
		return nil, connErr
	}
	return &RscsDB{
		sqliteDBFile: sqliteDBFile,
		db:           db}, nil
}

// DBFileName returns the db used.
func (r *RscsDB) DBFileName() string {
	return r.sqliteDBFile
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
	if key == "" {
		return 0, errors.New("insert empty key")
	}
	if len(key) > 255 {
		return 0, errors.New("key exceeds len")
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
	queryStr := fmt.Sprintf("DELETE FROM %s WHERE %s = $1",
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

// Update will give a row a new value.
func (r *RscsDB) Update(key, value string) (int, error) {
	if key == "" {
		return 0, errors.New("update empty key")
	}
	queryStr := fmt.Sprintf("UPDATE %s SET %s = $1 WHERE %s is $2",
		KVTableName, KVValueColumn, KVPrimaryKeyColumn)
	result, updateErr := r.db.Exec(queryStr, value, key)
	if updateErr != nil {
		return 0, updateErr
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
