package utils

import (
	"encoding/json"
	"net/http"
)

// Message generates a map of the response pieces
func Message(status bool, message string) (map[string]interface{}) {
	return map[string]interface{}{"status": status, "message": message}
}

// Respond prepares the encoded message to be sent in the response
func Respond(w http.ResponseWriter, data map[string] interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
