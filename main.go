package main

import (
	"net/http"
	"html/template"
	"./model"
	"./utils"
	"strconv"
	"fmt"
)

var dbInfo = "root:@/hermes"

var homeUrl string = "/"

var roomCreateUrl string = "/room/create"
var roomListUrl string = "/room"
var roomConversationUrl string = "/room/conversation"
var messageSendUrl string = "/room/conversation/send"
var signUpUrl string = "/user/signup"
var signInUrl string = "/user/signin"
var signOurUrl string = "/user/signout"

func main() {
	http.HandleFunc(homeUrl, HomeHandler)

	http.HandleFunc(roomCreateUrl, RoomCreateHandler)
	http.HandleFunc(roomListUrl, RoomListHandler)
	http.HandleFunc(roomConversationUrl, RoomConversationHandler)
	http.HandleFunc(messageSendUrl, MessageSendHandler)
	http.HandleFunc(signUpUrl, SignUpHandler)
	http.HandleFunc(signInUrl, SignInHandler)
	http.HandleFunc(signOurUrl, SignOutHandler)

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
	// TODO: Take note 2 db connections happening here
	var room *model.Room = model.GetRoomById(roomId)
	var conv []*model.Conversation = model.GetConversation(roomId)
	t.ExecuteTemplate(w, "room_detail", map[string]interface{}{
			"room": room,
			"conversation": conv,
		})
}

func MessageSendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var msg string = r.FormValue("content")
		roomId, _ := strconv.Atoi(r.FormValue("room_id"))
		userId, _ := strconv.Atoi(r.FormValue("user_id"))
		var message *model.Message = new(model.Message)
		message.Content = msg
		message.RoomId = roomId
		message.UserId = userId
		message.Insert()
	}
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var t = template.Must(template.New("SignUp").ParseFiles("templates/sign_up.html", "templates/header.html", "templates/footer.html"))
		t.ExecuteTemplate(w, "sign_up", nil)
	} else if r.Method == "POST" {
		var fname = r.FormValue("fname")
		var lname = r.FormValue("lname")
		var uname = r.FormValue("uname")
		var pword = r.FormValue("pword")
		var email = r.FormValue("email")

		var user *model.User = new(model.User)
		user.FirstName = fname
		user.LastName = lname
		user.Username = uname
		user.Password = pword
		user.Email = email
		if user.Insert() {
			fmt.Println("@" + uname + ": signed up")
			var response = `{"message": "Success"}`
			// Setting the response as json format
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(response))
		}
	}
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var t = template.Must(template.New("SignIn").ParseFiles("templates/sign_in.html", "templates/header.html", "templates/footer.html"))
		t.ExecuteTemplate(w, "sign_in", nil)
	} else if r.Method == "POST" {
		var uname = r.FormValue("uname")
		var pword = r.FormValue("pword")

		var user *model.User = new(model.User)
		user.Username = uname
		user.Password = pword

		var response []byte
		if user.Authenticate() {
			fmt.Println("@" + uname + ": signed in")
			response = []byte(`{"message": "Success"}`)
		} else {
			response = []byte(`{"message": "Mismatach username and password"}`)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}
}

func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var uname = r.FormValue("uname")
		fmt.Println("@" + uname + ": signed out")

		var response []byte = []byte(`{"message": "Success"}`)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

