package routes

import (
	"database/sql"
	"net/http"

	"github.com/dmarquinah/go_rooms/middlewares"
	"github.com/dmarquinah/go_rooms/services"
)

func createHostRoutes(mux *http.ServeMux, database *sql.DB) {
	mux.Handle("GET /host/self", middlewares.JWTmiddleware(handleGetHost(database)))
}

func handleGetHost(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services.GetHost(w, r, database)
	}
}
