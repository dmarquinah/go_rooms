package routes

import (
	"net/http"

	"github.com/dmarquinah/go_rooms/middlewares"
)

func createUserRoutes(mux *http.ServeMux) {
	mux.Handle("GET /user", middlewares.JWTmiddleware(handleGetUser))
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully obtained user"))
}
