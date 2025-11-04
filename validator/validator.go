package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// Validator provides common validation functions
type Validator struct {
	errors []ValidationError
}

// New creates a new validator instance
func New() *Validator {
	return &Validator{
		errors: make([]ValidationError, 0),
	}
}

// AddError adds a validation error
func (v *Validator) AddError(field, message string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// IsValid returns true if there are no validation errors
func (v *Validator) IsValid() bool {
	return len(v.errors) == 0
}

// Errors returns all validation errors
func (v *Validator) Errors() []ValidationError {
	return v.errors
}

// ErrorMessages returns all error messages as a single string
func (v *Validator) ErrorMessages() string {
	if v.IsValid() {
		return ""
	}
	var messages []string
	for _, err := range v.errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// Required checks if a string field is not empty
func (v *Validator) Required(field, value string) {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, "is required")
	}
}

// MinLength checks if a string meets minimum length requirement
func (v *Validator) MinLength(field, value string, min int) {
	if len(value) < min {
		v.AddError(field, fmt.Sprintf("must be at least %d characters", min))
	}
}

// MaxLength checks if a string does not exceed maximum length
func (v *Validator) MaxLength(field, value string, max int) {
	if len(value) > max {
		v.AddError(field, fmt.Sprintf("must not exceed %d characters", max))
	}
}

// Email validates email format
func (v *Validator) Email(field, value string) {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		v.AddError(field, "must be a valid email address")
	}
}

// Range checks if an integer is within a specified range
func (v *Validator) Range(field string, value, min, max int) {
	if value < min || value > max {
		v.AddError(field, fmt.Sprintf("must be between %d and %d", min, max))
	}
}

// Pattern checks if a string matches a regular expression pattern
func (v *Validator) Pattern(field, value, pattern, message string) {
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		v.AddError(field, "pattern validation failed")
		return
	}
	if !matched {
		v.AddError(field, message)
	}
}

// OneOf checks if a value is one of the allowed values
func (v *Validator) OneOf(field, value string, allowed []string) {
	for _, a := range allowed {
		if value == a {
			return
		}
	}
	v.AddError(field, fmt.Sprintf("must be one of: %s", strings.Join(allowed, ", ")))
}

// Custom allows adding custom validation logic
func (v *Validator) Custom(field string, valid bool, message string) {
	if !valid {
		v.AddError(field, message)
	}
}
