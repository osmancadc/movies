package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func GetRandomNumber(w http.ResponseWriter, r *http.Request) {

	response, err := http.Get("http://www.randomnumberapi.com/api/v1.0/random?min=100&max=1000&count=1")
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, err.Error())
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	number := string(body)
	SendResponse(w, http.StatusAccepted, number[1:len(number)-1])
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/random", GetRandomNumber).Methods("GET")
	router.HandleFunc("/users/register", CreateUser).Methods("POST")
	router.HandleFunc("/users/login", AccessUser).Methods("POST")
	router.HandleFunc("/movies/create", CreateMovie).Methods("POST")
	router.HandleFunc("/movies/delete", DeletePrivateMovie).Methods("DELETE")
	router.HandleFunc("/movies/delete/all", DeleteAllPrivateMovies).Methods("DELETE")
	router.HandleFunc("/movies/update/sales", UpdateSales).Methods("PATCH")
	router.HandleFunc("/movies/like", LikeMovie).Methods("POST")
	router.HandleFunc("/movies/get/public", GetPublicMovies).Methods("GET")
	router.HandleFunc("/movies/get/private", GetPrivateMovies).Methods("GET")
	router.HandleFunc("/movies/get/liked", GetLikedMovies).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", router))
}
