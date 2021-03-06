package user

import (
	"database/sql" // package to encode and decode the json into struct and vice versa
	"encoding/json"
	"fmt"
	"go-postgres-practice/models" // models package where User schema is defined
	"log"
	"net/http" // used to access the request and response object of the api
	"strconv"  // package used to covert string into int type

	"github.com/gorilla/mux" // used to get the params from the route

	_ "github.com/lib/pq" // postgres golang driver
)

// CreateUser creates a user in the postgres db
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		user     models.User
		insertID int64
		res      models.Response
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

// GetUser fetches a user by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	var (
		id   int
		err  error
		user models.User
	)

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get the userid from the request params, key is "id"
	params := mux.Vars(r)
	// convert the id type from str to int
	id, err = strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	user, err = getUser(int64(id))
	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	// send response
	json.NewEncoder(w).Encode(user)
}

// GetAllUser fetches all users
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	var (
		users []models.User
		err   error
	)

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	users, err = getAllUser()
	if err != nil {
		log.Fatalf("Unable to get all users. %v", err)
	}

	json.NewEncoder(w).Encode(users)
}

// UpdateUser updates a user's information in the DB
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var (
		params      map[string]string
		id          int
		err         error
		user        models.User
		updatedRows int64
		msg         string
		res         models.Response
	)

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params = mux.Vars(r)
	id, err = strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	updatedRows = updateUser(int64(id), user)
	msg = fmt.Sprintf("User updated successfully. Total rows/record affected %v", updatedRows)

	res.ID = int64(id)
	res.Message = msg

	json.NewEncoder(w).Encode(res)
}

// DeleteUser deletes a user record from the DB
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// TYPE ONLY declarations of variables
	var (
		params      map[string]string
		id          int
		err         error
		deletedRows int64
		msg         string
		res         models.Response
	)

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-", "*")
	w.Header().Set("Access-Control-Allow-", "DELETE")
	w.Header().Set("Access-Control-Allow-", "Content-Type")

	params = mux.Vars(r)
	id, err = strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to conver the string into int. %v", err)
	}

	deletedRows = deleteUser(int64(id))
	msg = fmt.Sprintf("User deleted successfully. Total rows/record affected %v", deletedRows)
	res.ID = int64(id)
	res.Message = msg

	json.NewEncoder(w).Encode(res)
}

// ============HANDLER FUNCTIONS============
func insertUser(user models.User) int64 {
	var (
		db           *sql.DB
		sqlStatement string
		id           int64
		err          error
	)

	db = models.CreateConnection()
	defer db.Close()

	// create insert SQL query; returning userid will return the newly created user's ID
	sqlStatement = `INSERT INTO users (name, location, age) VALUES ($1, $2, $3) RETURNING userid`

	err = db.QueryRow(sqlStatement, user.Name, user.Location, user.Age).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("inserted a single record %v", id)

	return id
}

func getUser(id int64) (models.User, error) {
	var (
		db           *sql.DB
		user         models.User
		sqlStatement string
		row          *sql.Row
		err          error
	)

	db = models.CreateConnection()
	defer db.Close()

	sqlStatement = `SELECT * FROM users WHERE userid=$1`
	row = db.QueryRow(sqlStatement, id)
	err = row.Scan(&user.ID, &user.Name, &user.Age, &user.Location)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return user, err
}

func getAllUser() ([]models.User, error) {
	var (
		db           *sql.DB
		users        []models.User
		sqlStatement string
		rows         *sql.Rows
		err          error
	)

	db = models.CreateConnection()
	defer db.Close()

	sqlStatement = `SELECT * FROM users`
	rows, err = db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Name, &user.Age, &user.Location)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		users = append(users, user)
	}

	return users, err
}

func updateUser(id int64, user models.User) int64 {
	var (
		db           *sql.DB
		sqlStatement string
		res          sql.Result
		err          error
		rowsAffected int64
	)

	db = models.CreateConnection()
	defer db.Close()
	sqlStatement = `UPDATE users SET name=$2, location=$3, age=$4 WHERE userid=$1`
	res, err = db.Exec(sqlStatement, id, user.Name, user.Location, user.Age)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the number of affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)
	return rowsAffected
}

func deleteUser(id int64) int64 {
	var (
		db           *sql.DB
		sqlStatement string
		res          sql.Result
		err          error
		rowsAffected int64
	)

	db = models.CreateConnection()
	defer db.Close()

	sqlStatement = `DELETE FROM users WHERE userid=$1`
	res, err = db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the number of affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)
	return rowsAffected
}
