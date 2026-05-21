package validator

import (
	"net/mail"
	"strings"
	"unicode/utf8"
)

type ValidationErrors map[string]string

func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

func Required(errors ValidationErrors, field, value string) {
	if strings.TrimSpace(value) == "" {
		errors[field] = field + " is required"
	}
}

func MinLength(errors ValidationErrors, field, value string, min int) {
	if utf8.RuneCountInString(value) < min {
		errors[field] = field + " must be at least " + string(rune(min+'0')) + " characters"
	}
}

func MaxLength(errors ValidationErrors, field, value string, max int) {
	if utf8.RuneCountInString(value) > max {
		errors[field] = field + " must not exceed " + string(rune(max+'0')) + " characters"
	}
}

func Email(errors ValidationErrors, field, value string) {
	if value == "" {
		return
	}
	_, err := mail.ParseAddress(value)
	if err != nil {
		errors[field] = field + " must be a valid email address"
	}
}

func MinValue(errors ValidationErrors, field string, value, min float64) {
	if value < min {
		errors[field] = field + " must be greater than or equal to minimum"
	}
}

func MaxValue(errors ValidationErrors, field string, value, max float64) {
	if value > max {
		errors[field] = field + " must be less than or equal to maximum"
	}
}

func InList(errors ValidationErrors, field, value string, allowed []string) {
	if value == "" {
		return
	}
	for _, a := range allowed {
		if value == a {
			return
		}
	}
	errors[field] = field + " must be one of: " + strings.Join(allowed, ", ")
}
