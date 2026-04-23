package v1

import (
	"net/http"

	"github.com/acakp/dumbchat/internal/adapter/postgres"
	"github.com/acakp/dumbchat/pkg/render"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RequireAdmin(dbpool *pgxpool.Pool, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("admin_session")
		if err != nil || cookie.Valid() != nil {
			render.Error(w, err, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if err = postgres.IsAdminSession(dbpool, cookie); err != nil {
			render.Error(w, err, http.StatusUnauthorized, "Unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}
