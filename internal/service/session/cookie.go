package session

import (
	"net/http"
	"time"
)

func SetTokenToCookie(w http.ResponseWriter, name, token string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Secure:   false, // Set to true in production with HTTPS
	}
	http.SetCookie(w, cookie)
}

func GetTokenFromCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func DeleteSessionCookie(w http.ResponseWriter, name string) {

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	})
}
