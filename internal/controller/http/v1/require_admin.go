package v1

import (
	"database/sql"
	"net/http"

	"github.com/acakp/dumbchat/internal/usecase"
)

func RequireAdmin(db *sql.DB, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("admin_session")
		if err != nil || cookie.Valid() != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if err = usecase.IsAdminSession(db, cookie); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
