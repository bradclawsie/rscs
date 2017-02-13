// Package server provides an Http interface to Rscs.
package server

import (
	"encoding/json"
	"net/http"
)

// StatusResult describes the system status.
type StatusResult struct {
	Alive  bool
	DBFile string
}

// Status returns the sytem status as JSON.
func (s *RscsServer) Status(w http.ResponseWriter, r *http.Request) {
	jsonBytes, jsonErr := json.Marshal(StatusResult{
		Alive:  true,
		DBFile: s.rscsDB.DBFileName()})
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(jsonBytes)
	return
}
