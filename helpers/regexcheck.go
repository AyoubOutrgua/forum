package helpers

import (
	"regexp"
)

func ValidateInfo(username, email, password string) bool {

	if len(email) > 50 || len(email) < 7 {
		return false
	} else if len(username) < 4 || len(username) > 15 {
		 return false
	} else if len(password) < 6 || len(password) > 20 {
		return false
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return !emailRegex.MatchString(email) 
		
	
	
	

}
