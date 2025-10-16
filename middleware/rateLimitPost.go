package middleware

import (
	"fmt"
	"time"

	"forum/database"
)

func RateLimitPost() bool {
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	userID := 2
	dateQuery := `
			SELECT creationDate
			FROM posts
			WHERE userId = ?
			ORDER BY creationDate DESC
			LIMIT 5;
			`

	dates := database.SelectDateOfLast5Posts(dateQuery, userID)
	for i, d := range dates {
		fmt.Println(i+1, " : ", d)
	}

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
		fmt.Println("LEN DATES:", len(dates))
		fmt.Println("DIF", dif)
		if len(dates) == 5 && dif <= 1 {
			return true
		}
	}

	// countQuery := `
	// 		SELECT count(*)
	// 		FROM posts
	// 		WHERE userId = ? AND creationDate BETWEEN ? AND ?;
	// 		`
	// numberOfPosts := database.CountPostsBetweenTimes(countQuery, userID, dateCreation, timeNow)
	// lastTime, err := time.Parse("2006-01-02 15:04:05", dateCreation)
	// timeNow2, err2 := time.Parse("2006-01-02 15:04:05", timeNow)
	// if err != nil || err2 != nil {
	// 	fmt.Println(err)
	// 	return true
	// }
	// diff := timeNow2.Sub(lastTime)
	// fmt.Println("LAST TIME:", lastTime.Format(time.RFC3339Nano))
	// fmt.Println("TIME NOW :", timeNow2.Format(time.RFC3339Nano))
	// fmt.Println("DIFF SEC :", diff.Minutes())
	// if diff.Seconds() <= 20 && numberOfPosts > 1 {
	// 	fmt.Println("Rak 3y9tiiiiiiiiiiiiiiiiiiiiiiiiiiiii Hbas chwyaaaaaaaaaaaaaaaaaaaaaaaa !!!!!!!!!!!!!!!!!!!!!!")
	// 	return true
	// }
	return false
}
