package core

import (
	"encoding/json"
	"net/http"
)

// ReadRequest parses application/json request bodies into value, or returns an error.
func ReadRequest(r *http.Request, value any) error {
	switch r.Header.Get("Content-Type") {
	case "application/json":
		err := json.NewDecoder(r.Body).Decode(value)
		if err != nil {
			return BadRequest("Failed to parse malformed request body, must be a valid JSON object")
		}

		return nil
	default:
		return BadRequest("Unsupported request Content-Type, must be application/json")
	}
}
