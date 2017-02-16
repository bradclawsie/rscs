// Package server provides an Http interface to Rscs.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// StatusResult describes the system status.
type StatusResult struct {
	Alive  bool
	DBFile string
	Uptime string
}

// Status returns the system status as JSON.
func (s *RscsServer) Status(w http.ResponseWriter, r *http.Request) {
	uptime := fmt.Sprintf("%v", time.Since(s.start))
	jsonBytes, jsonErr := json.Marshal(StatusResult{
		Alive:  true,
		DBFile: s.rscsDB.DBFileName(),
		Uptime: uptime})
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(jsonBytes)
	return
}
