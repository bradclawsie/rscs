package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Update changes the value of an existing key/value pair.
func (s *RscsServer) Update(w http.ResponseWriter, r *http.Request) {
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

	rowCount, updateErr := s.rscsDB.Update(key, *v.Value)
	if updateErr != nil {
		http.Error(w, updateErr.Error(), http.StatusInternalServerError)
		return
	}

	if rowCount == 0 {
		e := fmt.Sprintf("no key '%s' found", key)
		http.Error(w, e, http.StatusNotFound)
		return
	}
	if rowCount != 1 {
		http.Error(w, "bad number of rows updated", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
