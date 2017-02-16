// Package server provides an Http interface to Rscs.
package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// Update changes the value of an existing key/value pair.
func (s *RscsServer) Update(w http.ResponseWriter, r *http.Request) {
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

	body, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}

	var v valueVerify
	umErr := json.Unmarshal(body, &v)
	if umErr != nil || v.Value == nil {
		http.Error(w, "Value JSON malformed", http.StatusBadRequest)
		return
	}
	rowCount, updateErr := s.rscsDB.Update(key, *v.Value)
	if updateErr != nil {
		http.Error(w, updateErr.Error(), http.StatusBadRequest)
		return
	}
	if rowCount != 1 {
		http.Error(w, "bad number of rows updated", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
