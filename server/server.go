// Package server provides an Https interface to Rscs.
package server

import (
	"github.com/bradclawsie/rscs/db"
)

// RscsServer contains the state values for the underlying database instance
// and for https routing.
type RscsServer struct {
	rscsDB *db.RscsDB
}

// NewRscsServer initializes a new RscsServer instance.
func NewRscsServer(sqliteDBFile string, readOnly bool) (*RscsServer, error) {
	rscsDB, dbErr := db.NewRscsDB(sqliteDBFile, readOnly)
	if dbErr != nil {
		return nil, dbErr
	}
	return &RscsServer{rscsDB: rscsDB}, nil
}
