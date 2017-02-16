// Package server provides an Http interface to Rscs.
package server

import (
	"encoding/json"
	"net/http"
)

// Get retrieves the value for the key passed on the URL path.
func (s *RscsServer) Get(w http.ResponseWriter, r *http.Request) {
	key, keyErr := extractKeyContext(r)
	if keyErr != nil {
		http.Error(w, keyErr.Error(), http.StatusInternalServerError)
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

	jsonBytes, jsonErr := json.Marshal(Value{Value: value})
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(jsonBytes)
	return
}
