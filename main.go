package main

import (
	"log"
	"net/http"

	"github.com/dmarquinah/go_rooms/db"
	"github.com/joho/godotenv"
)

const SERVER_PORT string = ":5000"

func main() {
	//Handle load of environment params
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading the .env file")
	}

	// DB connection logic lays here
	db.HandleDBConnection()

	// Define basic route handlers to endpoints
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// types.Login(w, r, database)
	})
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		// types.HandleUser(w, r, database)
	})
	http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		// types.HandleTask(w, r, database)
	})

	println("Server listening on port " + SERVER_PORT)
	http.ListenAndServe(SERVER_PORT, nil)

}
