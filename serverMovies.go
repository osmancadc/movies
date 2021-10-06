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
	router.HandleFunc("/users/register", CreateUser).Methods("POST")
	router.HandleFunc("/users/login", AccessUser).Methods("POST")
	router.HandleFunc("/movies/get/public", getPublicMovies).Methods("GET")
	router.HandleFunc("/movies/get/private", getPrivateMovies).Methods("GET")
	router.HandleFunc("/movies/create", CreateMovie).Methods("POST")
	// router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	// router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router))
}
