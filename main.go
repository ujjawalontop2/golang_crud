package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// User struct represents the data model
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type CustomResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var db *sql.DB

func init() {
	// Connect to MySQL database
	var err error
	db, err = sql.Open("mysql", "root:2019354@tcp(localhost:3306)/myapi")
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := mux.NewRouter()

	// Define API routes
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users", updateUser).Methods("PATCH")

	// Start the server
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	// Fetch all users from the database
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []User

	// Iterate over the rows
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	response := CustomResponse{
		Status:  "Success",
		Message: "Details Fetched Successfully",
		Data:    users,
	}

	// Convert the users to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User

	// Decode the request body into a User struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the new user into the database
	result, err := db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", user.Name, user.Age)
	if err != nil {
		log.Fatal(err)
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	response := CustomResponse{
		Status:  "Success",
		Message: fmt.Sprintf("User ID %d created successfully", id),
		Data:    user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func updateUser(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	// Decode the request body into a User struct
	var updateData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var setClause string
	var values []interface{}
	for key, value := range updateData {
		switch key {
		case "name":
			setClause += "name = ?, "
			values = append(values, value.(string))
		case "age":
			setClause += "age = ?, "
			values = append(values, value.(int))
		}
	}

	// Remove the trailing comma from the setClause
	setClause = strings.TrimSuffix(setClause, ", ")

	// Construct and execute the UPDATE query
	_, err = db.Exec("UPDATE users SET "+setClause+" WHERE id = ?", append(values, id)...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := CustomResponse{
		Status:  "Success",
		Message: fmt.Sprintf("User ID %s updated successfully", id),
		Data:    updateData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
