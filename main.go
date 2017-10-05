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

type Response struct {
	Id interface{} `json:"id"`
}

func main() {
	databaseURI := os.Getenv("MYSQL_URL")

	db, err = sql.Open("mysql", databaseURI)
	check(err)
	defer db.Close()

	err = db.Ping()
	check(err)

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
	var responseId int
	response := Response{}

	r.ParseForm()
	for key := range r.PostForm {
		message = key
	}

	id := random(0, 99999)
	check(err)
	_, err = db.Exec("INSERT INTO messages(id, message) VALUES(?, ?)", id, message)
	check(err)

	err = db.QueryRow("SELECT id FROM messages WHERE id = ?", id).Scan(&responseId)
	check(err)
	response = Response{Id: responseId}
	js, _ := json.Marshal(response)
	check(err)
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "%s", js)
}

func getMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var message string
	err = db.QueryRow("SELECT message FROM messages WHERE id = ?", id).Scan(&message)
	check(err)

	fmt.Fprintf(w, "%v\n", message)
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
