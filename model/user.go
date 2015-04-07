package model

type User struct {
	Id			int		`db:"id"`
	Username	string	`db:"username"`
	FirstName	string	`db:"first_name"`
	LastName	string	`db:"last_name"`
	Password	string	`db:"password"`
	Email		string	`db:"email"`
}
