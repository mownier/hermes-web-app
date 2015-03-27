package main

import (
	"net/http"
	"html/template"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/room/create", CreateRoomHandler)
	http.HandleFunc("/room/create/process", ProcessCreateRoomHandler)
	http.HandleFunc("/room", RoomListHandler)
	http.ListenAndServe(":4321", nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var t = template.Must(template.New("Home").ParseFiles("templates/home.html", "templates/header.html", "templates/footer.html"))
	var array = []int{1, 2, 3}
	var info = map[string]interface{}{"numbers": array}
	t.ExecuteTemplate(w, "home", info);
}

func CreateRoomHandler(w http.ResponseWriter, r* http.Request) {
	var t = template.Must(template.New("CreateRoom").ParseFiles("templates/create_room.html", "templates/header.html", "templates/footer.html"))
	t.ExecuteTemplate(w, "create_room", nil)
}

func RoomListHandler(w http.ResponseWriter, r* http.Request) {
	var t = template.Must(template.New("RoomList").ParseFiles("templates/room_list.html", "templates/header.html", "templates/footer.html"))
	t.ExecuteTemplate(w, "room_list", nil)
}

func ProcessCreateRoomHandler(w http.ResponseWriter, r* http.Request) {
	var roomName string = r.FormValue("room_name")

	conn, err := sql.Open("mysql", "root:@/hermes")
	
	defer conn.Close()

	w.Write([]byte(roomName));
}

