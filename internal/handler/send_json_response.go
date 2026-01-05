package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// TODO : later change interface{}(any) to have only successResponse or errorResponse
func SendJSONResponse(w http.ResponseWriter, statusCode int, resp interface{}) {
	// Serialize Response
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(resp); err != nil {
		log.Printf("Failed to encode JSON: %v", err)

		// Fallback if encode failed
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Send Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(buf.Bytes())
}
