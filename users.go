package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/robbert229/jwt"
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
		SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	existUser, err := VerifyUser(register, newUser.Email)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if existUser {
		SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = InsertUser(newUser)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	SendResponse(w, http.StatusAccepted, "User created successfully")
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
		SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	existUser, err := VerifyUser(login, accessUser.Email, accessUser.Password)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !existUser {
		SendResponse(w, http.StatusBadRequest, "The user doesn't exists")
		return
	}

	token, err := GenerateToken(accessUser.Email)
	if err != nil {
		log.Printf("Error generando el token, %+v", err)
		SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := AccessReponse{
		Message: "User loged successfully",
		Token:   token,
	}

	SendResponse(w, http.StatusAccepted, response, true)
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

func GetUserID(token string) int {
	id := -1

	email, err := GetUserEmail(token)
	if err != nil {
		return id
	}

	db, err := GetDatabase()
	if err != nil {
		fmt.Printf("error with the database connection: %v", err)
		return id
	}

	defer db.Close()

	result, err := db.Query(`SELECT id FROM users where email = ?`, email)
	if err != nil {
		log.Println(err.Error())
	}
	result.Next()

	err = result.Scan(&id)
	if err != nil {
		return id
	}
	return id
}

func GetUserEmail(token string) (string, error) {
	algorithm := jwt.HmacSha256("SecretPassword")

	claims, err := algorithm.Decode(token)
	if err != nil {
		return "", err
	}

	email, err := claims.Get("email")
	if err != nil {
		return "", err
	}
	return email.(string), nil
}
