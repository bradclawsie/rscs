// Package server provides an Http interface to Rscs.
package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

// GetResult contains the value corresponding to a key.
type GetResult struct {
	Value string
}

// Get retrieves the value for the key passed on the URL path.
func (s *RscsServer) Get(w http.ResponseWriter, r *http.Request) {
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
