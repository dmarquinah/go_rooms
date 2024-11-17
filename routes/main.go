package routes

import (
	"database/sql"
	"net/http"
)

func BuildRouter(database *sql.DB) *http.ServeMux {
	router := http.NewServeMux()
	createAuthRoutes(router, database)
	createUserRoutes(router, database)
	createHostRoutes(router, database)

	return router
}
