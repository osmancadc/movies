package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	register = 1
	login    = 2
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a Valid Task Data")
	}

	json.Unmarshal(reqBody, &newUser)

	isValid, err := ValidateParameters(newUser)
	if !isValid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	existUser, err := VerifyUser(register, newUser.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	if existUser {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("The email is already registered")
		return
	}

	err = InsertUser(newUser)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode("User created successfully")
}

func AccessUser(w http.ResponseWriter, r *http.Request) {
	var accessUser User

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a Valid Task Data")
	}

	json.Unmarshal(reqBody, &accessUser)

	isValid, err := ValidateParameters(accessUser)
	if !isValid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	existUser, err := VerifyUser(login, accessUser.Email, accessUser.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	if !existUser {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("The user doesn't exists")
		return
	}

	token, err := GenerateToken(accessUser.Email)
	if err != nil {
		log.Printf("Error generando el token, %+v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	response := AccessReponse{
		Message: "User loged successfully",
		Token:   token,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonResponse)
}

func VerifyUser(verificationType int, parameters ...string) (bool, error) {
	db, err := GetDatabase()
	if err != nil {
		fmt.Printf("error with the database connection: %v", err)
		return false, err
	}

	defer db.Close()

	if verificationType == register {
		result, err := db.Query(`SELECT * FROM users where email like ?;`, parameters[0])
		if err != nil {
			log.Println(err.Error())
		}
		return result.Next(), nil
	} else {
		result, err := db.Query(`SELECT * FROM users where email = ? and password = ?;`, parameters[0], parameters[1])
		if err != nil {
			log.Println(err.Error())
		}
		return result.Next(), nil
	}
}

func InsertUser(user User) error {

	if user.Email == "" || user.Password == "" {
		return errors.New("invalid user")
	}

	db, err := GetDatabase()
	if err != nil {
		fmt.Printf("error with the database connection: %v", err)
		return err
	}

	defer db.Close()

	insertSentence, err := db.Prepare("INSERT INTO movies.users (name, email, password) VALUES(?, ?, ?);")
	if err != nil {
		return err
	}

	defer insertSentence.Close()

	_, err = insertSentence.Exec(user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}
