package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	log.Printf("WriteJSON - Status: %d, Headers: %+v", status, w.Header())
	return json.NewEncoder(w).Encode(data)

}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func WriteJSONError(w http.ResponseWriter, statusCode int, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(statusCode)
	log.Printf("WriteJSONError - Status: %d, Headers: %+v", statusCode, w.Header())

	switch msg := message.(type) {
	case string:
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": []map[string]string{
				{"error": msg},
			},
		})
	case map[string]string:
		var errorList []map[string]string
		for key, val := range msg {
			errorList = append(errorList, map[string]string{key: val})
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"errors": errorList})
	default:
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": []map[string]string{
				{"error": "An unexpected error occurred"},
			},
		})
	}
}
