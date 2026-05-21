package validator

import (
	"testing"
)

func TestRequired(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		hasErr  bool
	}{
		{"empty string", "", true},
		{"whitespace only", "   ", true},
		{"valid value", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidationErrors{}
			Required(errors, "field", tt.value)
			if tt.hasErr && !errors.HasErrors() {
				t.Error("Expected error but got none")
			}
			if !tt.hasErr && errors.HasErrors() {
				t.Error("Expected no error but got one")
			}
		})
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		hasErr bool
	}{
		{"valid email", "test@example.com", false},
		{"invalid email", "notanemail", true},
		{"empty email", "", false}, // empty is not validated by Email, use Required
		{"missing domain", "test@", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidationErrors{}
			Email(errors, "email", tt.value)
			if tt.hasErr && !errors.HasErrors() {
				t.Error("Expected error but got none")
			}
			if !tt.hasErr && errors.HasErrors() {
				t.Errorf("Expected no error but got: %v", errors)
			}
		})
	}
}

func TestMinValue(t *testing.T) {
	tests := []struct {
		name   string
		value  float64
		min    float64
		hasErr bool
	}{
		{"above min", 10, 5, false},
		{"equal to min", 5, 5, false},
		{"below min", 3, 5, true},
		{"zero below min", 0, 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidationErrors{}
			MinValue(errors, "field", tt.value, tt.min)
			if tt.hasErr && !errors.HasErrors() {
				t.Error("Expected error but got none")
			}
			if !tt.hasErr && errors.HasErrors() {
				t.Error("Expected no error but got one")
			}
		})
	}
}

func TestInList(t *testing.T) {
	allowed := []string{"income", "expense", "transfer"}

	tests := []struct {
		name   string
		value  string
		hasErr bool
	}{
		{"valid value", "income", false},
		{"another valid", "expense", false},
		{"invalid value", "other", true},
		{"empty value", "", false}, // empty is not validated by InList
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidationErrors{}
			InList(errors, "type", tt.value, allowed)
			if tt.hasErr && !errors.HasErrors() {
				t.Error("Expected error but got none")
			}
			if !tt.hasErr && errors.HasErrors() {
				t.Error("Expected no error but got one")
			}
		})
	}
}
