package model

type Message struct {
	Id 				int 	`db:"id"`
	Message 		string	`db:"string"`
	DateCreated 	string	`db:"date_created"`
	UserId 			int 	`db:"user_id"`
	RoomId 			int 	`db:"room_id"`
}
