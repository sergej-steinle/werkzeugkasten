package werkzeugkasten

import (
	"regexp"
	"testing"
)

func TestValidator_Valid(t *testing.T) {
	var v Validator
	if !v.Valid() {
		t.Error("expected Valid() to return true for empty FieldErrors")
	}

	v.AddFieldError("name", "must not be blank")
	if v.Valid() {
		t.Error("expected Valid() to return false after adding a field error")
	}
}

func TestValidator_Valid_WithNonFieldErrors(t *testing.T) {
	var v Validator
	if !v.Valid() {
		t.Error("expected Valid() to return true initially")
	}

	v.AddNonFieldError("something went wrong")
	if v.Valid() {
		t.Error("expected Valid() to return false after adding a non-field error")
	}
}

func TestValidator_AddNonFieldError(t *testing.T) {
	var v Validator
	v.AddNonFieldError("error one")
	v.AddNonFieldError("error two")
	v.AddNonFieldError("error one")

	if len(v.NonFieldErrors) != 3 {
		t.Errorf("expected 3 non-field errors but got %d", len(v.NonFieldErrors))
	}

	if v.NonFieldErrors[0] != "error one" {
		t.Errorf("expected first error to be %q, got %q", "error one", v.NonFieldErrors[0])
	}
	if v.NonFieldErrors[1] != "error two" {
		t.Errorf("expected second error to be %q, got %q", "error two", v.NonFieldErrors[1])
	}
}

func TestValidator_AddFieldError(t *testing.T) {
	var v Validator
	v.AddFieldError("email", "is required")
	v.AddFieldError("email", "must be valid")

	if len(v.FieldErrors) != 1 {
		t.Errorf("expected 1 field error but got %d", len(v.FieldErrors))
	}

	if v.FieldErrors["email"] != "is required" {
		t.Errorf("expected first error to be kept, got %q", v.FieldErrors["email"])
	}
}

func TestValidator_CheckField(t *testing.T) {
	var v Validator
	v.CheckField(true, "name", "must not be blank")
	if !v.Valid() {
		t.Error("expected no error when check passes")
	}

	v.CheckField(false, "name", "must not be blank")
	if v.Valid() {
		t.Error("expected error when check fails")
	}
	if v.FieldErrors["name"] != "must not be blank" {
		t.Errorf("expected error message %q, got %q", "must not be blank", v.FieldErrors["name"])
	}
}

func TestNotBlank(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"non-empty string", "hello", true},
		{"empty string", "", false},
		{"only spaces", "   ", false},
		{"tabs", "\t\t", false},
		{"text with spaces", "  hello  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NotBlank(tt.value); got != tt.expected {
				t.Errorf("NotBlank(%q) = %v, want %v", tt.value, got, tt.expected)
			}
		})
	}
}

func TestMaxChars(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		n        int
		expected bool
	}{
		{"under limit", "ab", 3, true},
		{"at limit", "abc", 3, true},
		{"over limit", "abcd", 3, false},
		{"empty string", "", 3, true},
		{"unicode chars", "äöü", 3, true},
		{"unicode over limit", "äöüß", 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxChars(tt.value, tt.n); got != tt.expected {
				t.Errorf("MaxChars(%q, %d) = %v, want %v", tt.value, tt.n, got, tt.expected)
			}
		})
	}
}

func TestMinChars(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		n        int
		expected bool
	}{
		{"over minimum", "abcd", 3, true},
		{"at minimum", "abc", 3, true},
		{"under minimum", "ab", 3, false},
		{"empty string", "", 1, false},
		{"zero minimum", "", 0, true},
		{"unicode chars at minimum", "äöü", 3, true},
		{"unicode chars under minimum", "äö", 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinChars(tt.value, tt.n); got != tt.expected {
				t.Errorf("MinChars(%q, %d) = %v, want %v", tt.value, tt.n, got, tt.expected)
			}
		})
	}
}

func TestPermittedValue(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		allowed  []string
		expected bool
	}{
		{"contained", "a", []string{"a", "b", "c"}, true},
		{"not contained", "d", []string{"a", "b", "c"}, false},
		{"empty allowed", "a", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PermittedValue(tt.value, tt.allowed...); got != tt.expected {
				t.Errorf("PermittedValue(%q) = %v, want %v", tt.value, got, tt.expected)
			}
		})
	}
}

func TestPermittedValue_Int(t *testing.T) {
	if !PermittedValue(1, 1, 2, 3) {
		t.Error("expected PermittedValue(1, 1,2,3) to be true")
	}
	if PermittedValue(4, 1, 2, 3) {
		t.Error("expected PermittedValue(4, 1,2,3) to be false")
	}
}

func TestMatches(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		rx       *regexp.Regexp
		expected bool
	}{
		{"digits only match", "12345", regexp.MustCompile(`^\d+$`), true},
		{"digits only no match", "123abc", regexp.MustCompile(`^\d+$`), false},
		{"lowercase match", "hello", regexp.MustCompile(`^[a-z]+$`), true},
		{"lowercase no match", "Hello", regexp.MustCompile(`^[a-z]+$`), false},
		{"empty string", "", regexp.MustCompile(`^.+$`), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Matches(tt.value, tt.rx); got != tt.expected {
				t.Errorf("Matches(%q) = %v, want %v", tt.value, got, tt.expected)
			}
		})
	}
}

func TestIsEmail(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"valid simple", "user@example.com", true},
		{"valid with dots", "first.last@example.com", true},
		{"valid with plus", "user+tag@example.com", true},
		{"valid with subdomain", "user@mail.example.com", true},
		{"valid with hyphen domain", "user@my-domain.com", true},
		{"missing @", "userexample.com", false},
		{"missing domain", "user@", false},
		{"missing local part", "@example.com", false},
		{"missing tld", "user@example", false},
		{"single char tld", "user@example.c", false},
		{"empty string", "", false},
		{"spaces", "user @example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmail(tt.value); got != tt.expected {
				t.Errorf("IsEmail(%q) = %v, want %v", tt.value, got, tt.expected)
			}
		})
	}
}
