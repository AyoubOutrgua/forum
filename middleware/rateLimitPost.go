package middleware

import (
	"fmt"
	"net/http"
	"time"

	"forum/helpers"
)

func RateLimitPost(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timeNow := time.Now().Format("2006-01-02 15:04:05")

		cookie, errSession := r.Cookie("session")
		if errSession != nil || cookie.Value == "" {
			fmt.Println("Error session makaynach 1")
		}
		cookieID := cookie.Value

		userID := helpers.GetUserID(cookieID)

		dates := helpers.GetLastPostDates(userID)
		firstDate := ""
		if len(dates) > 1 {
			firstDate = dates[len(dates)-1]
		}

		if firstDate != "" {
			firstDateTime, err1 := time.Parse("2006-01-02 15:04:05", firstDate)
			dateNowTime, err2 := time.Parse("2006-01-02 15:04:05", timeNow)
			if err1 != nil || err2 != nil {
			}
			dif := dateNowTime.Sub(firstDateTime).Minutes()
			if len(dates) == 5 && dif <= 1 {
				http.Error(w, "Ktbti posts bzaf !!!!!!!!", http.StatusBadRequest)
			}
		}
		next.ServeHTTP(w, r)
	}
}
