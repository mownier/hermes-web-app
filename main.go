package main

import (
	"net/http"
	"html/template"
	"./model"
	"./utils"
)

var dbInfo = "root:@/hermes"

var homeUrl string = "/"

var roomCreateUrl string = "/room/create"
var roomListUrl string = "/room"
var roomConversationUrl string = "/room/conversation"

func main() {
	http.HandleFunc(homeUrl, HomeHandler)

	http.HandleFunc(roomCreateUrl, RoomCreateHandler)
	http.HandleFunc(roomListUrl, RoomListHandler)
	http.HandleFunc(roomConversationUrl, RoomConversationHandler)

	// Mandatory root-based resources
    serveSingle("/scripts/jquery-1.11.2.min.js", "./scripts/jquery-1.11.2.min.js")

	http.ListenAndServe(":4321", nil)
}

func serveSingle(pattern string, filename string) {
    http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, filename)
    })
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
		room := new(model.Room)
		room.Name = roomName
		hasError := room.Insert()
		if hasError == false {
			println(utils.AppendString("created room: '", room.Name, "'"))
		}
		http.Redirect(w, r, roomListUrl, http.StatusMovedPermanently)
	}
}

func RoomListHandler(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap {
		"add": func(a, b int) int {
				return a + b
			},
	}
	var t = template.Must(template.New("RoomList").Funcs(funcMap).ParseFiles("templates/room_list.html", "templates/header.html", "templates/footer.html"))
	
	var data = map[string]interface{}{
		"rooms": model.GetAllRooms(),
	}

	t.ExecuteTemplate(w, "room_list", data)
}

func RoomConversationHandler(w http.ResponseWriter, r *http.Request) {
	var roomId string = r.URL.Query().Get("id")
	var t = template.Must(template.New("RoomDetail").ParseFiles("templates/room_detail.html", "templates/header.html", "templates/footer.html"))
	var room *model.Room = model.GetRoomById(roomId)
	var data map[string]string = nil
	if room != nil {
		data = map[string]string {
			"room_name": room.Name,
		}
	}
	t.ExecuteTemplate(w, "room_detail", data)
}
