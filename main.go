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

type response struct {
	ID interface{} `json:"id"`
}

func main() {
	databaseURI := os.Getenv("MYSQL_URL")

	db, err := sql.Open("mysql", databaseURI)
	if err != nil {
		println("Error:", err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		println("Error:", err.Error())
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

	res, err := db.Exec("INSERT INTO messages(message) VALUES(?)", message)
	if err != nil {
		println("Error:", err.Error())
	} else {
		responseID, err := res.LastInsertId()
		if err != nil {
			println("Error:", err.Error())
		} else {
			response := response{ID: responseID}
			js, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "%s\n", js)
		}
	}
}

func getMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var message string
	err := db.QueryRow("SELECT message FROM messages WHERE id = ?", id).Scan(&message)
	if err != nil {
		println("Error:", err.Error())
	}

	fmt.Fprintf(w, "%v\n", message)
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
