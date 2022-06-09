package main

import (
	"database/sql"
	"encoding/json"

	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	id         int    `json:"id"`
	User_id    int    `json:"user_id"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Biin       string `json:"biin"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Passwrd    string `json:"passwrd"`
}
type JsonResponse struct {
	Type    string `json:"type"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "1234"
	DB_NAME     = "test"
)

func main() {
	// Init the mux router
	router := mux.NewRouter()

	// Route handles & endpoints

	// Get user by biin
	router.HandleFunc("/user/{biin}", GetUser2).Methods("GET")

	// Create a user
	router.HandleFunc("/user", CreateUser).Methods("POST")

	// serve the app
	//PerformPostJsonRequest()
	fmt.Println("Server at 8000")
	log.Fatal(http.ListenAndServe(":8000", router))

}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func DB() *sql.DB {
	connStr := "user=postgres password=1234 dbname=Test sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	//defer db.Close()
	return db
}

/*func PerformPostJsonRequest() {
	const myurl = "https://petstore.swagger.io/v2/user"
	requestBody := strings.NewReader(`
		{
			"userId": 111,
		    "biin": "960825300817",
  			"firstName": "TestName",
  			"lastName": "TestSurname",
  			"password": "12345",
  			"phone": "77077663456",
  			"email": "kmf.test@kz"
		}
	`)
	response, err := http.Post(myurl, "application/json", requestBody)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	content, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(content))

}*/

func GetUser(w http.ResponseWriter, r *http.Request) {
	//db := setupDB()
	db := DB()
	// Get all movies from movies table that don't have movieID = "1"
	rows, err := db.Query("SELECT * FROM s_users")

	// check errors
	checkErr(err)

	// var response []JsonResponse
	users := make([]User, 0)
	defer rows.Close()

	// Foreach user

	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.id, &user.Phone, &user.Biin, &user.First_name, &user.Last_name, &user.Email, &user.Passwrd, &user.User_id)

		// check errors
		checkErr(err)
		//users = append(users, user)
		users = append(users, user)
		/*User{
			User_id:    1,
			First_name: "Jk",
			Last_name:  "Jk",
			Biin:       "biin",
			Email:      "email",
			Phone:      "phone123",
			Passwrd:    "pas",
		})*/
		return
	}
	params := mux.Vars(r)
	for _, user := range users {
		if user.Biin != params["biin"] {
			fmt.Fprintf(w, "%d %d %s %s %s %s %s %s \n ", user.id, user.User_id, user.Biin, user.Email, user.First_name, user.Last_name, user.Phone, user.Passwrd)
		}
	}

	var response = JsonResponse{Type: "success", Success: true, Message: ""}

	json.NewEncoder(w).Encode(response)

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	User_ID := r.FormValue("user_id")
	first_name := r.FormValue("first_name")

	var response = JsonResponse{}

	if User_ID == "" || first_name == "" {
		response = JsonResponse{Type: "error", Message: "You are missing  parameter."}
		json.NewEncoder(w).Encode(response)
	}

	db := DB()

	fmt.Println("Inserting user with ID: " + User_ID + " and name: " + first_name)

	//var lastInsertID int
	result, err := db.Exec("INSERT INTO s_users(USER_ID, FIRST_NAME) VALUES($1, $2)", User_ID, first_name)
	defer db.Close()
	// check errors
	checkErr(err)
	if result != nil {
		fmt.Println(result)
	} // количество добавленных строк

	response = JsonResponse{Type: "success", Success: true, Message: "The user has been inserted successfully!"}

	json.NewEncoder(w).Encode(response)
}

/*
func (u *User) FindByID(ID string) (User, error) {
	if ID == "" {
		return User{}, errors.New("ID can not be empty")
	}
	_, err := strconv.Atoi(ID)
	if err != nil {
		return User{}, errors.New("ID must be a number")
	}
	db := DB()
	user := User{}
	if err != nil {
		return user, err
	}
	defer db.Close()
	rows, err := db.Query(
		`SELECT * FROM users WHERE biin = ` + u.Biin)
	if err != nil {
		return user, err
	}
	if rows.Next() {
		rows.Scan(&user.id, &user.Phone, &user.Biin, &user.First_name, &user.Last_name, &user.Email, &user.Passwrd, &user.User_id)
		return user, nil
	}
	return User{}, errors.New("No data found")
}
*/
func GetUser2(w http.ResponseWriter, r *http.Request) {
	db := DB()
	sqlStatement := `SELECT * FROM s_users WHERE id=$1;`
	var user User
	row := db.QueryRow(sqlStatement, 3)
	err := row.Scan(&user.id, &user.Phone, &user.Biin, &user.First_name, &user.Last_name, &user.Email, &user.Passwrd, &user.User_id)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return
	case nil:
		fmt.Println(user)
	default:
		panic(err)
	}
}

/*func GetCar(w http.ResponseWriter, r *http.Request) {
	db := DB()

	params := mux.Vars(r)

	var user User

	db.First(&user , params["biin"])

	json.NewEncoder(w).Encode(&user)

  }*/
