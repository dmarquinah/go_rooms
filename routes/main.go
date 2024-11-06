package routes

import "net/http"

func BuildMux() *http.ServeMux {
	mux := http.NewServeMux()
	createAuthRoutes(mux)
	createUserRoutes(mux)

	return mux
}
