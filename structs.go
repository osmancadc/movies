package main

type Movie struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Duration int    `json:"duration"`
	Director string `json:"director"`
}

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccessReponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
