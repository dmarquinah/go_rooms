package routes

import (
	"database/sql"
	"net/http"

	"github.com/dmarquinah/go_rooms/middlewares"
	"github.com/dmarquinah/go_rooms/types"
)

func createUserRoutes(mux *http.ServeMux, database *sql.DB) {
	mux.Handle("GET /user/self", middlewares.JWTmiddleware(handleGetUser(database)))
}

func handleGetUser(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		types.GetUser(w, r, database)
	}
}
