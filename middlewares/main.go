package middlewares

import (
	"net/http"

	"github.com/dmarquinah/go_rooms/crypto"
)

func JWTmiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := crypto.GetJWTFromRequest(w, r)
		if token != nil {
			valid := crypto.ValidateJWT(*token, "id")

			if !valid {
				http.Error(w, "Invalid Auth token", http.StatusUnauthorized)
				return
			}

			next(w, r)
		}
	})
}
