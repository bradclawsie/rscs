// Package server provides an Https interface to Rscs.
package server

import (
	"encoding/json"
	"errors"
	"github.com/bradclawsie/rscs/db"
	"net/http"
	"strings"
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
	w.Write(jsonBytes)
	return
}

// GetResult contains the value corresponding to a key.
type GetResult struct {
	Value string
}

// Get retrieves the value for the key passed on the URL path.
func (s *RscsServer) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET supported", http.StatusMethodNotAllowed)
		return
	}
	pathElts := strings.Split(r.URL.Path, "/")
	if len(pathElts) == 0 {
		http.Error(w, "bad request path", http.StatusBadRequest)
		return
	}
	key := pathElts[len(pathElts)-1]
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}
	value, found, getErr := s.rscsDB.Get(key)
	if getErr != nil {
		http.Error(w, getErr.Error(), http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "no value found", http.StatusNotFound)
		return
	}
	jsonBytes, jsonErr := json.Marshal(GetResult{Value: value})
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(jsonBytes)
	return
}
