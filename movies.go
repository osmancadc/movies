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

func GetPublicMovies(w http.ResponseWriter, r *http.Request) {
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

func GetPrivateMovies(w http.ResponseWriter, r *http.Request) {
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

func GetLikedMovies(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	if !ValidateToken(token) {
		SendResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	user := GetUserID(token)

	db, err := GetDatabase()

	if err != nil {
		fmt.Printf("error with the database connection: %v", err)
		SendResponse(w, http.StatusInternalServerError, err.Error())
	}

	defer db.Close()

	results, err := db.Query(`SELECT m.id,title,duration,director,gender,premiere_year,sales  from moviesXlikes mx 
								inner join users u  on mx.user = u.id 
								inner join movies m on mx.movie = m.id
								where u.id = ?`, user)
	if err != nil {
		log.Println(err.Error())
	}
	defer results.Close()

	moviesLiked, err := ExtractMovies(results)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, err.Error())
	}

	SendResponse(w, http.StatusAccepted, moviesLiked, true)
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

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	var newMovie Movie

	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	if !ValidateToken(token) {
		SendResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	newMovie.Director = GetUserID(token)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		SendResponse(w, http.StatusBadRequest, "Invalid movie data")
		return
	}

	json.Unmarshal(reqBody, &newMovie)

	if newMovie.Name == "" || newMovie.Duration == 0 || newMovie.PremiereYear == "" {
		SendResponse(w, http.StatusBadRequest, "Invalid data to create a movie")
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

func DeleteAllPrivateMovies(w http.ResponseWriter, r *http.Request) {

	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	if !ValidateToken(token) {
		SendResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	director := GetUserID(token)

	db, err := GetDatabase()
	if err != nil {
		fmt.Printf("error with the database connection: %v", err)
		SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer db.Close()

	_, err = db.Query("DELETE FROM movies WHERE director = ? and public = 0", director)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	SendResponse(w, http.StatusAccepted, "All movies were deleted successfully")
}

func DeletePrivateMovie(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	if !ValidateToken(token) {
		SendResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	var deleteMovie Movie
	deleteMovie.Director = GetUserID(token)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		SendResponse(w, http.StatusBadRequest, "Invalid movie data")
		return
	}

	json.Unmarshal(reqBody, &deleteMovie)

	if deleteMovie.ID <= 0 {
		SendResponse(w, http.StatusBadRequest, "Invalid id to delete")
		return
	}

	db, err := GetDatabase()
	if err != nil {
		fmt.Printf("error with the database connection: %v", err)
		SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer db.Close()

	deleteSentence, err := db.Prepare("DELETE FROM movies WHERE director = ? and id = ? and public = 0")
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer deleteSentence.Close()

	_, err = deleteSentence.Exec(deleteMovie.Director, deleteMovie.ID)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, err.Error())
	}

	SendResponse(w, http.StatusAccepted, "Movie deleted successfully")
}

func UpdateSales(w http.ResponseWriter, r *http.Request) {

	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	if !ValidateToken(token) {
		SendResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	var updateMovie Movie
	updateMovie.Director = GetUserID(token)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		SendResponse(w, http.StatusBadRequest, "Invalid movie data")
		return
	}

	json.Unmarshal(reqBody, &updateMovie)

	if updateMovie.Sales <= 0 || updateMovie.ID <= 0 {
		SendResponse(w, http.StatusBadRequest, "Invalid data to update")
		return
	}

	db, err := GetDatabase()
	if err != nil {
		fmt.Printf("error with the database connection: %v", err)
		SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer db.Close()

	updateSentence, err := db.Prepare("UPDATE movies.movies SET sales=? WHERE director = ? and id = ?")
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer updateSentence.Close()

	_, err = updateSentence.Exec(updateMovie.Sales, updateMovie.Director, updateMovie.ID)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, err.Error())
	}

	SendResponse(w, http.StatusAccepted, "Movie's sales updated successfully")
}

func LikeMovie(w http.ResponseWriter, r *http.Request) {

	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	if !ValidateToken(token) {
		SendResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	var likedMovie Movie
	user := GetUserID(token)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		SendResponse(w, http.StatusBadRequest, "Invalid movie")
		return
	}

	json.Unmarshal(reqBody, &likedMovie)

	if likedMovie.ID <= 0 {
		SendResponse(w, http.StatusBadRequest, "Invalid movie to like")
		return
	}

	db, err := GetDatabase()
	if err != nil {
		fmt.Printf("error with the database connection: %v", err)
		SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer db.Close()

	insertSentence, err := db.Prepare("INSERT INTO movies.moviesXlikes (movie, `user`) VALUES(?, ?)")
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer insertSentence.Close()

	_, err = insertSentence.Exec(likedMovie.ID, user)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, err.Error())
	}

	SendResponse(w, http.StatusAccepted, "Movie liked successfully")
}
