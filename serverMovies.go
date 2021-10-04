package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func indexRoute(w http.ResponseWriter, r *http.Request) {
	log.Println("API Created successfully")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/users", createUser).Methods("POST")
	// router.HandleFunc("/tasks", getTasks).Methods("GET")
	// router.HandleFunc("/tasks/{id}", getOneTask).Methods("GET")
	// router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	// router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router))
}
