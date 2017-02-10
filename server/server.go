// Package server provides an Https interface to Rscs.
package server

import (
	"encoding/json"
	"errors"
	"github.com/bradclawsie/rscs/db"
	"io"
	"net/http"
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

// SHA256Response encodes the db SHA256.
type SHA256Response struct {
	SHA256 string
}

// SHA256 sends a JSON response with the current DB's SHA256.
func (s *RscsServer) SHA256(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET supported", http.StatusMethodNotAllowed)
		return
	}
	jsonBytes, jsonErr := json.Marshal(SHA256Response{SHA256: s.rscsDB.SHA256()})
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	io.WriteString(w, string(jsonBytes))
	return
}
