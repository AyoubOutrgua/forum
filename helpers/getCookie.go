package helpers

import "net/http"

func GetCookieValue(w http.ResponseWriter, r *http.Request) string {
	cookie, errSession := r.Cookie("session")
	if errSession != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return ""
	}
	return cookie.Value
}
