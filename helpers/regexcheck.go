package helpers

import (
	"regexp"
)

func ValidateInfo(username, email, password string) string {
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	if !usernameRegex.MatchString(username) {
		return ("invalid username: must be 3-20 characters (letters, digits, underscore)")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ("invalid email format")
	}

	passwordRegex := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$`)
	if !passwordRegex.MatchString(password) {
		return "invalid password: must be at least 8 chars with upper, lower, and number"
	}
	return ""
}
