package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

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
		fmt.Println(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/messages/", postMessageHandler).Methods("POST")
	router.HandleFunc("/messages/{id:[0-9]+}", getMessageHandler).Methods("GET")
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

	id := random(0, 999999, w, r)
	fmt.Fprintf(w, "Returned Id: %v\n", id)

	_, err := db.Exec("INSERT INTO messages(id, message) VALUES(?, ?)", id, message)
	if err != nil {
		fmt.Fprintf(w, "Error: %v\n", err)
	} else {
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
		fmt.Fprintf(w, "Error: %v\n", err)
	}

	fmt.Fprintf(w, "%v\n", message)
}

func random(min int, max int, w http.ResponseWriter, r *http.Request) int {
	var returnedID string

	rand.Seed(time.Now().Unix())
	id := rand.Intn(max-min) + min
	fmt.Fprintf(w, "%v\n", id)

	err := db.QueryRow("SELECT * FROM messages WHERE id = ?", id).Scan(&returnedID)
	if err != nil {
		fmt.Fprintf(w, "Error: %v\n", err)
		random(0, 999999, w, r)
	}

	fmt.Fprintf(w, "Id being returned: %v\n", id)
	return id
}
