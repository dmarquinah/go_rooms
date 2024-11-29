package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	driver "database/sql/driver"

	"github.com/go-sql-driver/mysql"
)

const DB_HOST_KEY = "DB_HOST"
const DB_PORT_KEY = "DB_PORT"
const DB_PASSWORD_KEY = "DB_PASSWORD"
const DB_USER_KEY = "DB_USER"
const DB_NAME_KEY = "DB_NAME"

const defaultCollation = "utf8mb4_general_ci"

func HandleDBConnection() (*sql.DB, error) {
	connector, err := getDBConnector()
	if err != nil {
		// Panic so it ends app execution
		panic(err.Error())
	}
	// Open connection to Database
	database := sql.OpenDB(connector)

	database.SetConnMaxLifetime(time.Minute * 5)
	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(10)

	// Not necessary to close the connection
	// https://pkg.go.dev/database/sql#Open
	//defer database.Close()

	// Check connection status
	err = database.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Database connected successfully!")
	return database, nil
}

func getDBConnector() (driver.Connector, error) {
	dbHost := os.Getenv(DB_HOST_KEY)
	dbPort := os.Getenv(DB_PORT_KEY)
	dbUser := os.Getenv(DB_USER_KEY)
	dbPassword := os.Getenv(DB_PASSWORD_KEY)
	dbName := os.Getenv(DB_NAME_KEY)

	dbAddr := dbHost + ":" + dbPort

	cfg := mysql.Config{
		User:                 dbUser,
		Passwd:               dbPassword,
		Net:                  "tcp",
		Addr:                 dbAddr,
		DBName:               dbName,
		Loc:                  time.Local,
		ParseTime:            true,
		Collation:            defaultCollation,
		AllowNativePasswords: true,
	}

	return mysql.NewConnector(&cfg)
}
