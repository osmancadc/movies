package main

import (
	"log"
	"net/http"
	"strings"
)

func getMovies(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	if ValidateToken(token) {
		log.Println("User accesed successfully")
	}
}
