package utils

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Parse UUID from string
func ParseUUID(id string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, ErrInvalidUUID
	}
	return parsedUUID, nil
}

// ParseDate converts a string (YYYY-MM-DD) to time.Time
func ParseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}

	parsedTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, ErrInvalidDate
	}
	return parsedTime, nil
}

// ParseDatePtr converts a string (YYYY-MM-DD) to *time.Time
func ParseDatePtr(dateStr *string) (*time.Time, error) {
	if dateStr == nil || *dateStr == "" {
		return nil, nil
	}
	parsedTime, err := time.Parse("2006-01-02", *dateStr)
	if err != nil {
		return nil, ErrInvalidDate
	}
	return &parsedTime, nil
}

// Returns nil if the string is empty, otherwise returns the string pointer
func NilIfEmpty(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}

// Converts an empty string to NULL for SQL compatibility
func NullTimeIfEmpty(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

// Checks if only allowed parameters are used in the request.
func ValidateParams(r *http.Request, allowedParams []string) error {
	allowed := make(map[string]bool)
	for _, param := range allowedParams {
		allowed[param] = true
	}

	for param := range r.URL.Query() {
		if !allowed[param] {
			return fmt.Errorf("%w: %s", ErrInvalidParameter, param)
		}
	}
	return nil
}
