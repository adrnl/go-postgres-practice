package product

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go-postgres-practice/models"
)

// CreateProduct creates a product entry in the DB
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		product  models.Product
		insertID int64
		res      models.Response
	)

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Fatalf("Unable to decode the json request body. %v", err)
	}

	insertID = insertProduct(product)
	res.ID = insertID
	res.Message = "Product created successfully"

	json.NewEncoder(w).Encode(res)
}

// ============HANDLER FUNCTIONS============
func insertProduct(product models.Product) int64 {
	var (
		db           *sql.DB
		sqlStatement string
		id           int64
		err          error
	)

	db = models.CreateConnection()
	defer db.Close()

	sqlStatement = `INSERT INTO products (name, msrp) VALUES ($1, $2) RETURNING productid`

	err = db.QueryRow(sqlStatement, product.Name, product.MSRP).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	return id
}
