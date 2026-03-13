package core

import (
	"encoding/json"
	"net/http"
)

// WriteResponse writes a JSON response, or 204 if the value is nil
// to indicate there is no content.
func WriteResponse(w http.ResponseWriter, value interface{}) {
	if value == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
