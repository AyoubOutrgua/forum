package helpers

import (
	"regexp"
)

func ValidateInfo(username, email, password string) string {
	ErrorMessage := ""
	if len(email) > 50 || len(email) < 7 {
		ErrorMessage = "Email must be between 7 and 50 characters"
	} else if len(username) < 4 || len(username) > 15 {
		ErrorMessage = "Username must be at least 3 characters to 14 characters"
	} else if len(password) < 6 || len(password) > 20 {
		ErrorMessage = "Password must be at least 6 characters to 14 characters"
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ("invalid email format")
	}

	return ErrorMessage
}
