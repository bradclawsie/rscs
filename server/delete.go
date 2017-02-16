// Package server provides an Http interface to Rscs.
package server

import (
	"net/http"
)

// Delete removes a row identified by a key.
func (s *RscsServer) Delete(w http.ResponseWriter, r *http.Request) {
	key, keyErr := extractKeyContext(r)
	if keyErr != nil {
		http.Error(w, keyErr.Error(), http.StatusInternalServerError)
		return
	}

	rowCount, deleteErr := s.rscsDB.Delete(key)
	if deleteErr != nil {
		http.Error(w, deleteErr.Error(), http.StatusInternalServerError)
		return
	}
	if rowCount != 1 {
		http.Error(w, "bad number of rows updated", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
