package core

import (
	"encoding/json"
	"mime"
	"net/http"
)

// ReadRequest parses application/json request bodies into value, or returns an error.
func ReadRequest(r *http.Request, value any) error {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		return BadRequest("Unsupported request Content-Type, must be application/json")
	}

	if err := json.NewDecoder(r.Body).Decode(value); err != nil {
		return BadRequest("Failed to parse malformed request body, must be a valid JSON object")
	}

	return nil
}
