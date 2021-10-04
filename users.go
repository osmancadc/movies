package main

import (
	"log"
	"net/http"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	// var newUser User

	// reqBody, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	fmt.Fprintf(w, "Insert a Valid Task Data")
	// }

	// json.Unmarshal(reqBody, &newTask)
	// newTask.ID = len(tasks) + 1
	// tasks = append(tasks, newTask)

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(newTask)
	log.Println("This is creation of users")

}

// func funcionPrueba() {
// 	log.Println("Funcion prueba")
// }
