package api

import (
	"net/http"

	"github.com/totoledao/auction-house/internal/jsonutils"
)

func (api *Api) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if !api.Sessions.Exists(r.Context(), "AuthenticatedUserId") {
				_ = jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{
					"message": "Unauthorized",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
}
