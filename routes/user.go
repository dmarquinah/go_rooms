package routes

import (
	"database/sql"
	"net/http"

	"github.com/dmarquinah/go_rooms/middlewares"
	"github.com/dmarquinah/go_rooms/services"
)

func createUserRoutes(mux *http.ServeMux, database *sql.DB) {
	mux.Handle("GET /user/self", middlewares.JWTmiddleware(handleGetUser(database)))
}

func handleGetUser(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services.GetUser(w, r, database)
	}
}
