// Package server provides an Http interface to Rscs.
package server

import (
	"net/http"
	"strings"
)

// Delete removes a row identified by a key.
func (s *RscsServer) Delete(w http.ResponseWriter, r *http.Request) {
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
