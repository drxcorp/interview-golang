package utils

import (
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}

	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func ValidatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}
	return true
}

func SanitizeString(input string) string {
	input = strings.ReplaceAll(input, "'", "''")
	input = strings.ReplaceAll(input, ";", "")
	input = strings.ReplaceAll(input, "--", "")
	return input
}

func IsNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func TruncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen]
	}
	return s
}
