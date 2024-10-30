package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const DB_HOST_KEY = "DB_HOST"
const DB_PORT_KEY = "DB_PORT"
const DB_PASSWORD_KEY = "DB_PASSWORD"
const DB_USER_KEY = "DB_USER"
const DB_NAME_KEY = "DB_NAME"

func HandleDBConnection() {
	// Open connection to Database
	database, err := sql.Open("mysql", getDSN())

	if err != nil {
		// Panic so it ends app execution
		panic(err.Error())
	}
	defer database.Close()

	// Check connection status
	err = database.Ping()
	if err != nil {
		panic(err.Error())
	}
	log.Println("Database connected successfully!")
}

func getDSN() string {
	dbHost := os.Getenv(DB_HOST_KEY)
	dbPort := os.Getenv(DB_PORT_KEY)
	dbUser := os.Getenv(DB_USER_KEY)
	dbPassword := os.Getenv(DB_PASSWORD_KEY)
	dbName := os.Getenv(DB_NAME_KEY)
	return dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName
}
