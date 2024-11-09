package routes

import (
	"database/sql"
	"net/http"

	"github.com/dmarquinah/go_rooms/types"
)

func createAuthRoutes(mux *http.ServeMux, database *sql.DB) {
	mux.HandleFunc("POST /auth/login", handlePostLogin(database))
	mux.HandleFunc("POST /auth/register", handlePostRegister(database))
}

func handlePostLogin(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		types.LoginUser(w, r, database)
	}
}

func handlePostRegister(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		types.RegisterUser(w, r, database)
	}
}
