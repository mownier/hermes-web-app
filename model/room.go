package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"fmt"
)

type Room struct {
	Id				int		`db:"id"`
	Name			string	`db:"name"`
	DateCreated 	string	`db:"date_created"`
	DateModified	string	`db:"date_modified"`
	Deleted 		int 	`db:"deleted"`
}

func (r *Room) Insert() bool {
	var valid bool
	db, _ := sql.Open("mysql", dbInfo)
	date := time.Now()
	_, err := db.Exec("INSERT INTO rooms (name, date_created, date_modified) VALUES (?, ?, ?)", r.Name, date, date)
	if err != nil {
		// TODO : handle errors
		fmt.Print("error: ")
		fmt.Println(err)
		valid = false
	} else {
		valid = true
	}
	defer db.Close()
	return valid
}

func (r *Room) Update() bool {
	var valid bool
	db, _ := sql.Open("mysql", dbInfo)
	date := time.Now()
	_, err := db.Exec("UPDATE rooms SET name=?, date_modified=? WHERE id=?", r.Name, date, r.Id)
	if err != nil {
		// TODO : handle errors
		fmt.Print("error: ")
		fmt.Println(err)
		valid = false
	} else {
		valid = true
	}
	defer db.Close()
	return valid
}

func (r *Room) Delete() bool {
	var valid bool
	db, _ := sql.Open("mysql", dbInfo)
	date := time.Now()
	_, err := db.Exec("UPDATE rooms SET deleted=?, date_modified=? WHERE id=?", 1, date, r.Id)
	if err != nil {
		// TODO : handle errors
		fmt.Print("error: ")
		fmt.Println(err)
		valid = false
	} else {
		valid = true
	}
	defer db.Close()
	return valid
}

func GetAllRooms() []*Room {
	db, _ := sql.Open("mysql", dbInfo)
	rows, _ := db.Query("SELECT * FROM rooms WHERE deleted=?", 0)
	var rooms []*Room = toRoomArray(rows)
	defer db.Close()
	return rooms
}

func GetRoomById(id string) *Room {
	db, _ := sql.Open("mysql", dbInfo)
	rows, err := db.Query("SELECT * FROM rooms WHERE id=?", id)
	if err != nil {
		// TODO : handle errors
		fmt.Print("error: ")
		fmt.Println(err)

		return nil
	} else {
		rooms := toRoomArray(rows)
		return rooms[0]
	}
}

func toRoomArray(rows *sql.Rows) []*Room {
	var data []*Room = make([]*Room, 0)
	for rows.Next() {
		var room *Room = new(Room)
		err := rows.Scan(&room.Id, &room.Name, &room.DateCreated, &room.DateModified, &room.Deleted)
		// TODO : handle errors
		if err != nil {
			fmt.Print("error: ")
			fmt.Println(err)
		}

		data = append(data, room)
	}
	return data
}
