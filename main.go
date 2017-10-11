package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

type response struct {
	ID interface{} `json:"id"`
}

func main() {
	databaseURI := os.Getenv("MYSQL_URL")

	db, err = sql.Open("mysql", databaseURI)
	if err != nil {
		fmt.Println("An error occured while opening the database - Error details:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("An error occurred while testing the connection with the database - Error details:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/messages/", postMessageHandler).Methods("POST")
	router.HandleFunc("/messages/{id:[0-9]+}", getMessageHandler).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(notFound)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func postMessageHandler(w http.ResponseWriter, r *http.Request) {
	var message string

	r.ParseForm()
	for key := range r.PostForm {
		message = key
	}

	res, err := db.Exec("INSERT INTO messages(message) VALUES(?)", message)
	if err != nil {
		fmt.Fprintf(w, "An error occurred while trying to insert the message into the database - Error details: %v\n", err)
	} else {
		id, err := res.LastInsertId()
		if err != nil {
			fmt.Fprintf(w, "An error occured while trying to retrieve the id of the last inserted entry - Error details: %v\n", err)
		}
		jsonResponse := response{ID: id}
		js, _ := json.Marshal(jsonResponse)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s\n", js)
	}
}

func getMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var message string
	err := db.QueryRow("SELECT message FROM messages WHERE id = ?", id).Scan(&message)
	if err != nil {
		fmt.Fprintf(w, "An error occured while trying to retrieve a message from the database - Error details: %v\n", err)
	}

	fmt.Fprintf(w, "%v\n", message)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Unfortunately you have requested an endpoint that doesn't exist, please refer to the relevant documentation.")
}
