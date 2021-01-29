package middleware

import (
	"database/sql" // package to encode and decode the json into struct and vice versa
	"fmt"
	// models package where User schema is defined
	"log" // used to access the request and response object of the api
	"os"  // used to read the environment variable
	// package used to covert string into int type
	// used to get the params from the route

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	var (
		err error
		db  *sql.DB
	)

	// load .env file
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env")
	}

	// open connection to database
	db, err = sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	// check connection to database
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successful Connection")

	return db
}
