package routes

import (
	"database/sql"
	"net/http"

	"github.com/dmarquinah/go_rooms/services"
)

func createAuthRoutes(mux *http.ServeMux, database *sql.DB) {
	mux.HandleFunc("POST /auth/login", handlePostLogin(database))
	mux.HandleFunc("POST /auth/register", handlePostRegister(database))
	mux.HandleFunc("POST /auth/host/login", handlePostHostLogin(database))
	mux.HandleFunc("POST /auth/host/register", handlePostHostRegister(database))
}

func handlePostLogin(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services.LoginUser(w, r, database)
	}
}

func handlePostRegister(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services.RegisterUser(w, r, database)
	}
}

func handlePostHostLogin(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services.LoginHost(w, r, database)
	}
}

func handlePostHostRegister(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services.RegisterHost(w, r, database)
	}
}
