package middleware

import (
	"database/sql" // package to encode and decode the json into struct and vice versa
	"encoding/json"
	"fmt"
	"net/http"
	// models package where User schema is defined
	"log" // used to access the request and response object of the api
	"os"  // used to read the environment variable
	// package used to covert string into int type
	// used to get the params from the route

	"github.com/adrnl/go-postgres-practice/models"
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

// CreateUser creates a user in the postgres db
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		user     models.User
		insertID int64
		res      response
	)

	// set the header to content type x-www-form-urlencoded
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	// allow all origin to handle cors issue
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// decode json into user object
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatalf("Unable to decode the json request body. %v", err)
	}

	insertID = insertUser(user)
	res.ID = insertID
	res.Message = "User created successfully"

	json.NewEncoder(w).Encode(res)
}

// HANDLER FUNCTIONS
func insertUser(user models.User) int64 {
	var (
		db           *sql.DB
		sqlStatement string
		id           int64
		err          error
	)

	db = createConnection()
	defer db.Close()

	// cfreate insert SQL query; returning userid will return the newly created user's ID
	sqlStatement = `INSERT INTO users (name, location, age) VALUES ($1, $2, $3) RETURNING userid`

	err = db.QueryRow(sqlStatement, user.Name, user.Location, user.Age).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("inserted a single record %v", id)

	return id
}
