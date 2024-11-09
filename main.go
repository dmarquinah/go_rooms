package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dmarquinah/go_rooms/db"
	"github.com/dmarquinah/go_rooms/routes"
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
	database := db.HandleDBConnection()

	mux := routes.BuildMux(database)

	fmt.Println("Server listening on port " + SERVER_PORT)
	log.Fatal(http.ListenAndServe(SERVER_PORT, mux))
}
