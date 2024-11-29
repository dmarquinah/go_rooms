package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dmarquinah/go_rooms/database/cache"
	"github.com/dmarquinah/go_rooms/database/db"
	"github.com/dmarquinah/go_rooms/routes"
	"github.com/joho/godotenv"
)

const DEFAULT_SERVER_PORT string = ":5000"

func main() {
	//Handle load of environment params
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading the .env file")
	}

	// Initialize Redis
	err = cache.InitRedis()
	if err != nil {
		log.Fatalf("Failed to connect to Cache: %v", err)
		os.Exit(1)
	}

	// DB connection logic lays here
	database, err := db.HandleDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to Database: %v", err)
		os.Exit(1)
	}

	router := routes.BuildRouter(database)

	updatedRouter := routes.SetupGlobalMiddlewares(router)

	SERVER_PORT := os.Getenv("PORT")
	if SERVER_PORT == "" {
		SERVER_PORT = DEFAULT_SERVER_PORT
	}

	fmt.Println("Server listening on port " + SERVER_PORT)
	if err := http.ListenAndServe(SERVER_PORT, updatedRouter); err != nil {
		log.Fatalf("error on init server")
		os.Exit(1)
	}
}
