package middlewares

import (
	"context"
	"net/http"

	"github.com/dmarquinah/go_rooms/crypto"
	"github.com/dmarquinah/go_rooms/types"
	"github.com/dmarquinah/go_rooms/utils"
)

func JWTmiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := crypto.GetJWTFromRequest(w, r)
		if token != nil {
			valid := crypto.ValidateJWT(*token)

			if !valid {
				types.WriteErrorResponse(w, "Invalid Auth token", http.StatusUnauthorized)
				return
			}

			id, err := crypto.GetFieldFromJWT(*token, "id")
			if err != nil {
				types.WriteErrorResponse(w, "Invalid Auth token", http.StatusUnauthorized)
				return
			}

			role, err := crypto.GetFieldFromJWT(*token, "role")
			if err != nil {
				types.WriteErrorResponse(w, "Invalid Auth token", http.StatusUnauthorized)
				return
			}

			// Setting up the ID from validated user into the context so it can be used down the call chain
			ctx := context.WithValue(r.Context(), utils.IdKey, id)
			ctx_roles := context.WithValue(ctx, utils.RoleKey, role)
			r = r.WithContext(ctx_roles)

			next.ServeHTTP(w, r)
		}
	})
}
