package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Insert adds a new key/value pair.
func (s *RscsServer) Insert(w http.ResponseWriter, r *http.Request) {
	key, keyErr := extractKeyContext(r)
	if keyErr != nil {
		http.Error(w, keyErr.Error(), http.StatusInternalServerError)
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
