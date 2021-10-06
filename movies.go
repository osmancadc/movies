package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	private = 1
	public  = 2
)

func getPublicMovies(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	if !ValidateToken(token) {
		SendResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	movies, err := GetMoviesFromDatabase(public)
	if err != nil {
		log.Println(err.Error())
		SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	SendResponse(w, http.StatusAccepted, movies, true)
}

func getPrivateMovies(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	if !ValidateToken(token) {
		SendResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	email, err := GetUserEmail(token)
	if err != nil {
		log.Println(err.Error())
		SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	movies, err := GetMoviesFromDatabase(private, email)
	if err != nil {
		log.Println(err.Error())
		SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	SendResponse(w, http.StatusAccepted, movies, true)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	var newMovie Movie

	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	newMovie.Director = GetUserID(token)

	if !ValidateToken(token) {
		SendResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		SendResponse(w, http.StatusBadRequest, "Invalid movie data 1")
		return
	}

	json.Unmarshal(reqBody, &newMovie)

	log.Printf("Movie: %+v", newMovie)
	if newMovie.Name == "" || newMovie.Duration == 0 || newMovie.PremiereYear == "" {
		SendResponse(w, http.StatusBadRequest, "Invalid movie data 2")
		return
	}

	db, err := GetDatabase()
	if err != nil {
		fmt.Printf("error with the database connection: %v", err)
		SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer db.Close()

	insertSentence, err := db.Prepare("INSERT INTO movies.movies (title, duration, director, gender, premiere_year, sales, public) VALUES(?, ?, ?, ?, ?, ?, 0); ")
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer insertSentence.Close()

	_, err = insertSentence.Exec(newMovie.Name, newMovie.Duration, newMovie.Director, newMovie.Gender, newMovie.PremiereYear, newMovie.Sales)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, err.Error())
	}

	SendResponse(w, http.StatusAccepted, "Movie created successfully")
}

func GetMoviesFromDatabase(visibility int, parameters ...string) ([]Movie, error) {
	db, err := GetDatabase()

	if err != nil {
		fmt.Printf("error with the database connection: %v", err)
		return []Movie{}, err
	}

	defer db.Close()

	if visibility == public {

		results, err := db.Query(`SELECT id,title,duration,director,gender,premiere_year,sales FROM movies where public=true`)
		if err != nil {
			log.Println(err.Error())
		}
		defer results.Close()

		return ExtractMovies(results)

	} else {
		results, err := db.Query(`SELECT m.id,title,duration,director,gender,premiere_year,sales FROM movies m inner join users u on m.director = u.id where m.public is false and u.email = ?`, parameters[0])
		if err != nil {
			log.Println(err.Error())
		}
		defer results.Close()

		return ExtractMovies(results)
	}
}

func ExtractMovies(results *sql.Rows) ([]Movie, error) {
	movies := []Movie{}
	var movie Movie

	for results.Next() {
		err := results.Scan(&movie.ID, &movie.Name, &movie.Duration, &movie.Director, &movie.Gender, &movie.PremiereYear, &movie.Sales)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	return movies, nil
}
