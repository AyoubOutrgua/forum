package helpers

import (
	"regexp"
)
// ValidateInfo checks if the provided username, email, and password meet the required format and length constraints.
func ValidateInfo(username, email, password string) bool {

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if len(email) > 50 || len(email) < 7 {
		return false
	} else if len(username) < 4 || len(username) > 15 {
		 return false
	} else if len(password) < 6 || len(password) > 20 {
		return false
		}else if  !emailRegex.MatchString(email) {
			return false
		}
		

		
	
	return true
	

}
