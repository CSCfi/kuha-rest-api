package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)

}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(data); err != nil {
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("request body is required")
		}
		return err
	}

	return nil
}

func WriteJSONError(w http.ResponseWriter, statusCode int, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

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
