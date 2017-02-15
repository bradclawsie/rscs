// Package server provides an Http interface to Rscs.
package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// Insert adds a new key/value pair.
func (s *RscsServer) Insert(w http.ResponseWriter, r *http.Request) {
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

	type valueVerify struct {
		Value *string
	}
	var v valueVerify
	umErr := json.Unmarshal(body, &v)
	if umErr != nil || v.Value == nil {
		http.Error(w, "Value JSON malformed", http.StatusBadRequest)
		return
	}

	rowCount, insertErr := s.rscsDB.Insert(key, *v.Value)
	if insertErr != nil {
		http.Error(w, insertErr.Error(), http.StatusBadRequest)
		return
	}
	if rowCount != 1 {
		http.Error(w, "bad number of rows created", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}
