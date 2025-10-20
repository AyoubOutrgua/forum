package helpers

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type FlashData struct {
	Message  string
	Email    string
	Username string
}

func SetFlash(w http.ResponseWriter, message string, email string, username string) {
	data := FlashData{
		Message:  message,
		Email:    email,
		Username: username,
	}
	jsondata, _ := json.Marshal(data)
	encoded := url.QueryEscape(string(jsondata))

	http.SetCookie(w, &http.Cookie{
		Name:   "flash",
		Value:  string(encoded),
		Path:   "/",
		MaxAge: 60,
	})
}

func GetFlash(w http.ResponseWriter, r *http.Request) FlashData {
	var flash FlashData
	c, err := r.Cookie("flash")
	if err != nil {
		return flash
	}
	decoded, _ := url.QueryUnescape(c.Value)
	json.Unmarshal([]byte(decoded), &flash)

	http.SetCookie(w, &http.Cookie{
		Name:   "flash",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	return flash
}
