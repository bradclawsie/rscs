// Package server provides an Https interface to Rscs.
package server

import (
	"github.com/bradclawsie/rscs/db"
	"errors"
)

// RscsServer contains the state values for the underlying database instance
// and for https routing.
type RscsServer struct {
	rscsDB *db.RscsDB
}

// NewRscsServer initializes a new RscsServer instance.
func NewRscsServer(rscsDB *db.RscsDB) (*RscsServer, error) {
	if rscsDB == nil {
		return nil, errors.New("nil rscsDB")
	}
	return &RscsServer{rscsDB: rscsDB}, nil
}
