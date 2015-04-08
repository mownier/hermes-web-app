package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"fmt"
)

type Message struct {
	Id 				int 	`db:"id"`
	Content 		string	`db:"content"`
	DateCreated 	string	`db:"date_created"`
	DateModified	string	`db:"date_modified"`
	UserId 			int 	`db:"user_id"`
	RoomId 			int 	`db:"room_id"`
	Deleted 		int 	`db:"deleted"`
}

type Conversation struct {
	Room		*Room
	User 		*User
	Message 	*Message
}

func (m *Message) Insert() bool {
	var valid bool
	db, _ := sql.Open("mysql", dbInfo)
	date := time.Now()
	_, err := db.Exec("INSERT INTO messages (content, date_created, date_modified, user_id, room_id) VALUES (?, ?, ?, ?, ?)", m.Content, date, date, m.UserId, m.RoomId)
	if err != nil {
		// TODO: Handle errors
		fmt.Print("errorsr: ")
		fmt.Println(err)
		valid = false
	} else {
		valid = true
	}
	defer db.Close()
	return valid
}

func (m *Message) Update() bool {
	var valid bool
	db, _ := sql.Open("mysql", dbInfo)
	date := time.Now()
	_, err := db.Exec("UPDATE messages SET content=?, date_modified=? WHERE id=?", m.Content, date, m.Id)
	if err != nil {
		// TODO: Handle errors
		fmt.Print("error: ")
		fmt.Println(err)
		valid = false
	} else {
		valid = true
	}
	defer db.Close()
	return valid
}

func (m *Message) Delete() bool {
	var valid bool
	db, _ := sql.Open("mysql", dbInfo)
	date := time.Now()
	_, err := db.Exec("UPDATE messages SET date_modified=?, deleted=? WHERE id=?", date, 1, m.Id)
	if err != nil {
		// TODO: Handle errors
		fmt.Print("error: ")
		fmt.Println(err)
		valid = false
	} else {
		valid = true
	}
	defer db.Close()
	return valid
}

func GetConversation(roomId string) []*Conversation {
	db, _ := sql.Open("mysql", dbInfo)
	var query string = 
	`SELECT messages.*,
			users.username, users.first_name, users.last_name, users.email,
			rooms.name as room_name
		FROM messages
		INNER JOIN users ON messages.user_id=users.id
		INNER JOIN rooms ON messages.room_id=rooms.id
			WHERE messages.deleted=? AND rooms.id=?
		ORDER BY date_created DESC`
	rows, err := db.Query(query, 0, roomId)
	if err != nil {
		// TODO: Handle errors
		fmt.Print("error: ")
		fmt.Println(err)
	}
	convs := toConversationArray(rows)
	defer db.Close()
	return convs
}

func NewConversation() *Conversation {
	conversation := new(Conversation)
	conversation.Message = new(Message)
	conversation.Room = new(Room)
	conversation.User = new(User)
	return conversation
}

func toConversationArray(rows *sql.Rows) []*Conversation {
	var data []*Conversation = make([]*Conversation, 0)
	for rows.Next() {
		var conv *Conversation = NewConversation()
		err := rows.Scan(&conv.Message.Id, &conv.Message.Content, &conv.Message.UserId, &conv.Message.RoomId, &conv.Message.DateCreated, &conv.Message.DateModified, &conv.Message.Deleted,
						 &conv.User.Username, &conv.User.FirstName, &conv.User.LastName, &conv.User.Email,
						 &conv.Room.Name)
		conv.User.Id = conv.Message.UserId
		conv.Room.Id = conv.Message.RoomId

		// TODO : handle errors
		if err != nil {
			fmt.Print("error: ")
			fmt.Println(err)
		}

		data = append(data, conv)
	}
	return data
}

