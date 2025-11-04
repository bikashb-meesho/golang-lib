package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator_Required(t *testing.T) {
	v := New()
	
	v.Required("name", "")
	assert.False(t, v.IsValid())
	assert.Len(t, v.Errors(), 1)
	assert.Equal(t, "name", v.Errors()[0].Field)
	
	v = New()
	v.Required("name", "John")
	assert.True(t, v.IsValid())
}

func TestValidator_MinLength(t *testing.T) {
	v := New()
	
	v.MinLength("password", "123", 8)
	assert.False(t, v.IsValid())
	assert.Contains(t, v.ErrorMessages(), "must be at least 8 characters")
	
	v = New()
	v.MinLength("password", "12345678", 8)
	assert.True(t, v.IsValid())
}

func TestValidator_MaxLength(t *testing.T) {
	v := New()
	
	v.MaxLength("name", "a very long name that exceeds the limit", 10)
	assert.False(t, v.IsValid())
	
	v = New()
	v.MaxLength("name", "short", 10)
	assert.True(t, v.IsValid())
}

func TestValidator_Email(t *testing.T) {
	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{"valid email", "user@example.com", true},
		{"valid with subdomain", "user@mail.example.com", true},
		{"valid with plus", "user+tag@example.com", true},
		{"invalid - no @", "userexample.com", false},
		{"invalid - no domain", "user@", false},
		{"invalid - no tld", "user@example", false},
		{"invalid - spaces", "user @example.com", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			v.Email("email", tt.email)
			assert.Equal(t, tt.valid, v.IsValid())
		})
	}
}

func TestValidator_Range(t *testing.T) {
	v := New()
	
	v.Range("age", 5, 18, 65)
	assert.False(t, v.IsValid())
	
	v = New()
	v.Range("age", 100, 18, 65)
	assert.False(t, v.IsValid())
	
	v = New()
	v.Range("age", 30, 18, 65)
	assert.True(t, v.IsValid())
}

func TestValidator_OneOf(t *testing.T) {
	allowed := []string{"admin", "user", "guest"}
	
	v := New()
	v.OneOf("role", "invalid", allowed)
	assert.False(t, v.IsValid())
	assert.Contains(t, v.ErrorMessages(), "must be one of")
	
	v = New()
	v.OneOf("role", "admin", allowed)
	assert.True(t, v.IsValid())
}

func TestValidator_Pattern(t *testing.T) {
	v := New()
	
	// Test alphanumeric pattern
	v.Pattern("username", "user@123", `^[a-zA-Z0-9]+$`, "must be alphanumeric")
	assert.False(t, v.IsValid())
	
	v = New()
	v.Pattern("username", "user123", `^[a-zA-Z0-9]+$`, "must be alphanumeric")
	assert.True(t, v.IsValid())
}

func TestValidator_Custom(t *testing.T) {
	v := New()
	
	// Test custom validation
	isEven := func(n int) bool { return n%2 == 0 }
	
	v.Custom("number", isEven(5), "must be even")
	assert.False(t, v.IsValid())
	
	v = New()
	v.Custom("number", isEven(4), "must be even")
	assert.True(t, v.IsValid())
}

func TestValidator_MultipleErrors(t *testing.T) {
	v := New()
	
	v.Required("name", "")
	v.Email("email", "invalid")
	v.Range("age", 200, 1, 150)
	
	assert.False(t, v.IsValid())
	assert.Len(t, v.Errors(), 3)
	
	errorMsg := v.ErrorMessages()
	assert.Contains(t, errorMsg, "name")
	assert.Contains(t, errorMsg, "email")
	assert.Contains(t, errorMsg, "age")
}

func TestValidator_CompleteValidation(t *testing.T) {
	v := New()
	
	// Valid user data
	v.Required("name", "John Doe")
	v.MinLength("name", "John Doe", 2)
	v.Required("email", "john@example.com")
	v.Email("email", "john@example.com")
	v.Range("age", 30, 1, 150)
	v.OneOf("role", "user", []string{"admin", "user", "guest"})
	
	assert.True(t, v.IsValid())
	assert.Empty(t, v.ErrorMessages())
}

