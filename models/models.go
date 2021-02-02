package models

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
