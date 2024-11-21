package routes

import (
	"database/sql"
	"net/http"

	"github.com/dmarquinah/go_rooms/middlewares"
	"github.com/dmarquinah/go_rooms/types"
)

func createRoomRoutes(mux *http.ServeMux, database *sql.DB) {
	mux.Handle("GET /room/{id}", middlewares.JWTmiddleware(handleGetRoomById(database)))
	mux.Handle("POST /room", middlewares.JWTmiddleware(handlePostRoom(database)))
}

func handleGetRoomById(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		types.GetRoomById(w, r, database)
	}
}

func handlePostRoom(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		types.CreateRoom(w, r, database)
	}
}