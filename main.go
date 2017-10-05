package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func main() {
	//databaseURI := os.Getenv("MYSQL_URL")

	db, err = sql.Open("mysql", "mr5j21hizy137wdq:nlg9q2kkx99dqruc@irkm0xtlo2pcmvvz.chr7pe7iynqr.eu-west-1.rds.amazonaws.com:3306/nzwyicxlxzeuimx2") //databaseURI
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
	router.HandleFunc("/messages", postMessageHandler).Methods("POST")
	router.HandleFunc("messages/{id:[0-9]+}", getMessageHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func postMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	var message string
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &message)

	id := random(0, 99999)
	check(err)
	_, err := db.Exec("INSERT INTO messages(message, id) VALUES(?, ?)", message, id)
	check(err)
}

func getMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var message string
	err = db.QueryRow("SELECT message FROM messages WHERE id = ?", id).Scan(&message)
	check(err)

	fmt.Fprintf(w, "Message: %v\n", message)
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
