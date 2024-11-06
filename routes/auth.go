package routes

import "net/http"

func createAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/login", handlePostLogin)
}

func handlePostLogin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is a POST login"))
}
