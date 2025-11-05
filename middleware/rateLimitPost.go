package middleware

import (
	"net/http"
	"time"

	"forum/helpers"
)

func RateLimitPost(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, errSession := r.Cookie("session")
		if errSession != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		cookieID := cookie.Value

		userID := helpers.GetUserID(cookieID)

		dates := helpers.GetLastPostDates(userID)
		if len(dates) == 5 {
			firstDateTime, _ := time.Parse("2006-01-02 15:04:05", dates[len(dates)-1])
			dateNowTime := time.Now()

			if dateNowTime.Sub(firstDateTime).Minutes() <= 1 {
				helpers.Errorhandler(w, "Status Too Many Requests", http.StatusTooManyRequests)
				return
			}
		}

		next.ServeHTTP(w, r)
	}
}
