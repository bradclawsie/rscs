package server

import (
	"fmt"
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

	// Caller can ignore return value for deletes on keys not found.
	if rowCount == 0 {
		e := fmt.Sprintf("no key '%s' found", key)
		http.Error(w, e, http.StatusNotFound)
		return
	}
	if rowCount > 1 {
		http.Error(w, "bad number of rows updated", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
