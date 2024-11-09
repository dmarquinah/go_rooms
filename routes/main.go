package routes

import (
	"database/sql"
	"net/http"
)

func BuildMux(database *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	createAuthRoutes(mux, database)
	createUserRoutes(mux, database)

	return mux
}
