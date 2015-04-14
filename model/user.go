package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

type User struct {
	Id			int		`db:"id"`
	Username	string	`db:"username"`
	FirstName	string	`db:"first_name"`
	LastName	string	`db:"last_name"`
	Password	string	`db:"password"`
	Email		string	`db:"email"`
}

func (u *User) Insert() bool {
	var valid bool
	db, _ := sql.Open("mysql", dbInfo)
	var query = `
		INSERT INTO users (username, first_name, last_name, password, email)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := db.Exec(query, u.Username, u.FirstName, u.LastName, u.Password, u.Email)
	if err != nil {
		fmt.Print("error: ")
		fmt.Println(err)
		valid = false
	} else {
		valid = true
	}
	defer db.Close()

	return valid
}

func (u *User) Authenticate() bool {
	var valid bool
	db, _ := sql.Open("mysql", dbInfo)
	var query = `
		SELECT COUNT(*) FROM users WHERE username=? AND password=?
	`
	rows, err := db.Query(query, u.Username, u.Password)
	if err != nil {
		fmt.Print("error: ")
		fmt.Println(err)
		valid = false
	} else {
		var count int
		rows.Next()
		rows.Scan(&count)
		if count == 1 {
			valid = true
		} else {
			valid = false
		}
		
	}
	defer db.Close()

	return valid
}
