package usecase

import (
	"net/http"
)

func IssueAdminSession(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     "admin_session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   36000, // 10h
	}
	http.SetCookie(w, cookie)
}
