package utils

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// Helper function to convert `sql.NullString` to a regular string
func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)

}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func WriteJSONError(w http.ResponseWriter, status int, msg string) error {
	type envelope struct {
		Error string `json:"error"`
	}
	return WriteJSON(w, status, &envelope{Error: msg})
}
