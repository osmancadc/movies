package main

type Movie struct {
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"title"`
	Duration     int    `json:"duration" db:"duration"`
	Director     int    `json:"director" db:"director"`
	Gender       string `json:"gender" db:"gender"`
	PremiereYear string `json:"premiere_year" db:"premiere_year"`
	Sales        int    `json:"sales" db:"sales"`
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
