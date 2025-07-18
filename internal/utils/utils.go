package utils

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// query timeout duration
const QueryTimeout = 7 * time.Second

// Validator to be initialized once
var validate *validator.Validate
var once sync.Once

// Validator getter
func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())

		// Custom validation
		validate.RegisterValidation("key", func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			match, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, value)
			return match
		})
	})
	return validate
}

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

// Converts a string into sql.NullString with Valid=true
func NullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}

// Checks if only allowed parameters are used in the request.
func ValidateParams(r *http.Request, allowedParams []string) error {
	allowed := make(map[string]bool)
	for _, param := range allowedParams {
		allowed[param] = true
	}

	var invalidParams []string
	for param := range r.URL.Query() {
		if !allowed[param] {
			invalidParams = append(invalidParams, param)
		}
	}

	if len(invalidParams) > 0 {
		return fmt.Errorf("invalid query parameters: %s. Allowed parameters are: %s",
			strings.Join(invalidParams, ", "),
			strings.Join(allowedParams, ", "))
	}

	return nil
}

// Converts a string pointer to sql.NullString
func NullStringPtr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

// Converts a float64 pointer to sql.NullFloat64
func NullFloat64Ptr(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}

// Converts an int32 pointer to sql.NullInt32
func NullInt32Ptr(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}
