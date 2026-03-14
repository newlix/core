package core

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// WriteResponse writes a JSON response, or 204 if the value is nil
// to indicate there is no content.
func WriteResponse(w http.ResponseWriter, value any) {
	if value == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf.Bytes())
}
