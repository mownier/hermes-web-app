package main

import (
	"bytes"
	"net/http"
	"html/template"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var dbInfo = "root:@/hermes"

var roomCreateUrl = "/room/create"
var roomListUrl = "/room"

func main() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc(roomCreateUrl, RoomCreateHandler)
	http.HandleFunc(roomListUrl, RoomListHandler)

	http.ListenAndServe(":4321", nil)
}

func appendString(b ...string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(b); i++ {
		var s string = b[i]
		buffer.WriteString(s)
	}
	return buffer.String()
}

func convertToArray(rows *sql.Rows) []interface{} {
	var result []interface{} = make([]interface{}, 0)
	var id int
	var name string
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}

		var row map[string]interface{} = map[string]interface{}{
			"id": id,
			"name": name,
		}
		
		result = append(result, row)
		
	}
	return result
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var t = template.Must(template.New("Home").ParseFiles("templates/home.html", "templates/header.html", "templates/footer.html"))
	var array = []int{1, 2, 3}
	var info = map[string]interface{}{"numbers": array}
	t.ExecuteTemplate(w, "home", info);
}

func RoomCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var t = template.Must(template.New("CreateRoom").ParseFiles("templates/create_room.html", "templates/header.html", "templates/footer.html"))
		t.ExecuteTemplate(w, "create_room", nil)
	} else if r.Method == "POST" {
		var roomName string = r.FormValue("room_name")
		
		db, _ := sql.Open("mysql", dbInfo)
		_, err := db.Exec("INSERT INTO rooms (name) VALUES (?)", roomName)
		if err == nil {
			println(appendString("created room: '", roomName, "'"))
		}
		defer db.Close()

		http.Redirect(w, r, roomListUrl, http.StatusMovedPermanently)
	}
}

func RoomListHandler(w http.ResponseWriter, r *http.Request) {
	var t = template.Must(template.New("RoomList").ParseFiles("templates/room_list.html", "templates/header.html", "templates/footer.html"))
	
	// TODO: Get list of rooms from the database
	db, _ := sql.Open("mysql", dbInfo)
	rows, _ := db.Query("SELECT * FROM rooms")
	defer rows.Close()
	defer db.Close()

	var rooms = convertToArray(rows)
	var data = map[string]interface{}{
		"rooms": rooms,
	}
	t.ExecuteTemplate(w, "room_list", data)
}

