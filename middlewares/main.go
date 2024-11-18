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

			// Setting up the ID from validated user into the context so it can be used to further requests
			ctx := context.WithValue(r.Context(), utils.IdKey, id)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}
	})
}
