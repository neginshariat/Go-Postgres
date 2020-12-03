package middleware

import (
	"database/sql"
	"encoding/json"
	"first-go-postgres/models"
	"fmt"
	"log"
	"net/http"
	"strconv"

	//"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	/* 	err := godotenv.Load(".env")
	   	if err != nil {
	   		log.Fatal("Error loading .env file")
	   	} */
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/gouser?sslmode=disable")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("connected successfully")
	return db
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal("Unable to decode the request body. %v", err)
	}
	insertID := insertUser(user)

	res := response{
		ID:      insertID,
		Message: "User created successfully",
	}
	json.NewEncoder(w).Encode(res)
}

func insertUser(user models.User) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO users (name, location, age) VALUES ($1, $2, $3) RETURNING userid`

	var id int64
	//err := db.QueryRow(sqlStatement, user.Name, user.Location, user.Age).Scan(&id)
	err := db.QueryRow(sqlStatement, user.Name, user.Location, user.Age).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	fmt.Printf("inserted a single record %v", id)

	return id
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string to int. %v", err)
	}
	user, err := getUser(int64(id))
	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	json.NewEncoder(w).Encode(user)
}

func getUser(id int64) (models.User, error) {

	db := createConnection()
	defer db.Close()

	var user models.User

	sqlStatement := `SELECT * FROM users WHERE userid=$1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&user.ID, &user.Name, &user.Age, &user.Location)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}
	return user, err
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	users, err := getAllUsers()
	if err != nil {
		log.Fatalf("Unable to get all user.%v", err)
	}
	json.NewEncoder(w).Encode(users)
}

func getAllUsers() ([]models.User, error) {
	db := createConnection()
	defer db.Close()
	var users []models.User
	sqlStatement := `SELECT* FROM users`

	row, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query %v", err)
	}
	defer row.Close()

	for row.Next() {

		var user models.User

		err := row.Scan(&user.ID, &user.Name, &user.Age, &user.Location)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		users = append(users, user)
	}
	return users, err
}

func DeleteUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal("Unable to convert string to int", err)
	}

	deletedRows := deleteUser(int64(id))

	msg := fmt.Sprintf("User updated successfully", deletedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)

}

func deleteUser(id int64) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM users WHERE userid=$1`

	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatal("Unable to delete . %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Unable to check affected rows", err)
	}
	fmt.Println("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal("Unable to convert string to int. %v", err)
	}

	var user models.User

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal("unale to decode the request body", err)
	}

	updatedRows := updateUser(int64(id), user)

	msg := fmt.Sprintf(" User updated successfully . Total rows/record affected %v.", updatedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)

}

func updateUser(id int64, user models.User) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `UPDATE users SET name=$2, location=$3, age=$4 WHERE userid=$1`

	res, err := db.Exec(sqlStatement, id, user.Name, user.Location, user.Age)
	if err != nil {
		log.Fatal("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Unable to checking the affected rows", err)
	}
	fmt.Println("Total rows/records affected %v", rowsAffected)

	return rowsAffected
}
