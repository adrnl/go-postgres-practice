package product

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go-postgres-practice/models"

	"github.com/gorilla/mux"
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

//GetProduct fetches a product by ID
func GetProduct(w http.ResponseWriter, r *http.Request) {
	var (
		params  map[string]string
		id      int
		err     error
		product models.Product
	)

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params = mux.Vars(r)
	id, err = strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	product, err = getProduct(int64(id))
	if err != nil {
		log.Fatalf("Unable to get product. %v", err)
	}

	json.NewEncoder(w).Encode(product)
}

// GetAllProduct fetches all products
func GetAllProduct(w http.ResponseWriter, r *http.Request) {
	var (
		products []models.Product
		err      error
	)

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	products, err = getAllProduct()
	if err != nil {
		log.Fatalf("Unable to get all products. %v", err)
	}

	json.NewEncoder(w).Encode(products)
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

func getProduct(id int64) (models.Product, error) {
	var (
		db           *sql.DB
		product      models.Product
		sqlStatement string
		row          *sql.Row
		err          error
	)

	db = models.CreateConnection()
	defer db.Close()

	sqlStatement = `SELECT * FROM products WHERE productid=$1`
	row = db.QueryRow(sqlStatement, id)
	err = row.Scan(&product.ID, &product.Name, &product.MSRP)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned")
		return product, nil
	case nil:
		return product, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return product, err
}

func getAllProduct() ([]models.Product, error) {
	var (
		db           *sql.DB
		products     []models.Product
		sqlStatement string
		rows         *sql.Rows
		err          error
	)

	db = models.CreateConnection()
	defer db.Close()

	sqlStatement = `SELECT * FROM product`
	rows, err = db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.ID, &product.Name, &product.MSRP)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		products = append(products, product)
	}

	return products, err
}

func updateProduct(id int64, product models.Product) int64 {
	var (
		db           *sql.DB
		sqlStatement string
		res          sql.Result
		err          error
		rowsAffected int64
	)

	db = models.CreateConnection()
	defer db.Close()

	sqlStatement = `UPDATE product SET name=$2, msrp=$3 WHERE productid=$1`
	res, err = db.Exec(sqlStatement, id, product.Name, product.MSRP)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the number of affected rows. %v", err)
	}

	fmt.Printf("Total rows/records affected %v", rowsAffected)
	return rowsAffected
}

func deleteProduct(id int64) int64 {
	var (
		db           *sql.DB
		sqlStatement string
		res          sql.Result
		err          error
		rowsAffected int64
	)

	db = models.CreateConnection()
	defer db.Close()

	sqlStatement = `DELETE FROM products WHERE productid=$1`
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
