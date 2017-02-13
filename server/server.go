// Package server provides an Http interface to Rscs.
package server

import (
	"errors"
	"github.com/bradclawsie/rscs/db"
	"github.com/pressly/chi"
)

const (
	// KVRoute is the route for all key/val operations.
	KVRoute = "/v1/kv/:key"
	// StatusRoute is the route for system status.
	StatusRoute = "/v1/status"
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

// NewRouter provides a new chi router to pass to a server.
func (s *RscsServer) NewRouter() (*chi.Mux, error) {
	rtr := chi.NewRouter()

	rtr.Get(KVRoute, s.Get)
	rtr.Get(StatusRoute, s.Status)

	return rtr, nil
}
