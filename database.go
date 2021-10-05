package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func GetDatabase() (db *sql.DB, e error) {
	user := "root"
	password := "FcjrhrjNvyNsh4Jh"
	host := "tcp(34.139.20.237:3306)"
	database := "movies"
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", user, password, host, database))
	if err != nil {
		return nil, err
	}
	return db, nil
}
