package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" // package used to read the .env file
)

// Response Type
type Response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// User Schema
type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Age      int64  `json:"age"`
}

// Product Schema
type Product struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	MSRP int64  `json:"msrp"`
}

// CreateConnection establishes a connection to the DB
func CreateConnection() *sql.DB {
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
