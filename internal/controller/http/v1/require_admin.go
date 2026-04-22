package v1

import (
	"database/sql"
	"net/http"

	"github.com/acakp/dumbchat/internal/usecase"
	"github.com/acakp/dumbchat/pkg/render"
)

func RequireAdmin(db *sql.DB, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("admin_session")
		if err != nil || cookie.Valid() != nil {
			render.Error(w, err, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if err = usecase.IsAdminSession(db, cookie); err != nil {
			render.Error(w, err, http.StatusUnauthorized, "Unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}
