package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

func main() {
	databaseURI := os.Getenv("MYPOSTGRES_URL")

	db, err = sql.Open("postgres", databaseURI)
	check(err)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := mux.NewRouter()
	router.HandleFunc("/", Index)
	router.HandleFunc("/messages", PostMessageHandler).Methods("POST")
	router.HandleFunc("messages/{id:[0-9]+}", GetMessageHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func Index(w http.ResponseWriter, r *http.Request) {

}

func PostMessageHandler(w http.ResponseWriter, r *http.Request) {

}

func GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
