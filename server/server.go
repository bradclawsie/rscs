// Package server provides an Http interface to Rscs.
package server

import (
	"errors"
	"github.com/bradclawsie/rscs/db"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"time"
)

const (
	// KVRoutePrefix is the prefix for the kv route.
	KVRoutePrefix = "/v1/kv"
	// KVRoute is the route for all key/val operations.
	KVRoute = KVRoutePrefix + "/:key"
	// StatusRoute is the route for system status.
	StatusRoute = "/v1/status"
	// keyName is the string for 'key'.
	keyName = "key"
)

// Value corresponds to a row value.
type Value struct {
	Value string
}

// valueVerify is like Value but nil-able for some internal validation purposes.
type valueVerify struct {
	Value *string
}

// RscsServer contains the state values for the underlying database instance
// and for https routing.
type RscsServer struct {
	rscsDB *db.RscsDB
	start  time.Time
}

// NewRscsServer initializes a new RscsServer instance.
func NewRscsServer(rscsDB *db.RscsDB) (*RscsServer, error) {
	if rscsDB == nil {
		return nil, errors.New("nil rscsDB")
	}
	return &RscsServer{rscsDB: rscsDB, start: time.Now()}, nil
}

// NewRouter provides a new chi router to pass to a server.
func (s *RscsServer) NewRouter() (*chi.Mux, error) {
	rtr := chi.NewRouter()
	rtr.Use(middleware.Recoverer)

	rtr.Route(KVRoute, func(rtr chi.Router) {
		rtr.Use(insertKeyContext)
		rtr.Get("/", s.Get)
		rtr.Post("/", s.Insert)
		rtr.Put("/", s.Update)
		rtr.Delete("/", s.Delete)
	})

	rtr.Get(StatusRoute, s.Status)

	return rtr, nil
}
