package handler

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/akihiro21/goChat/handlers/database"
)

type Msg struct {
	Message string
}

var (
	tokens    []string
	msg       = Msg{}
	db        *sql.DB
	userDB    = database.NewUserDB()
	roomDB    = database.NewRoomDB()
	MessageDB = database.NewMessageDB()
	templates = make(map[string]*template.Template)
)

func Init(mux *http.ServeMux) {
	go H.run()
	db = database.ConnectDB()
	defer db.Close()
	SessionInit()
	templatesInit()
	HandlerInit(mux)
}

func templatesInit() {
	templates["login"] = loadTemplate("login")
	templates["create"] = loadTemplate("create")
	templates["room"] = loadTemplate("room")
	templates["chat"] = loadTemplate("chat")
	templates["adminChat"] = loadTemplate("adminChat")
	templates["admin"] = loadTemplate("admin")
	templates["user"] = loadTemplate("user")
}

func HandlerInit(mux *http.ServeMux) {
	mux.HandleFunc("/", root)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/create", create)
	mux.HandleFunc("/room", room)
	mux.HandleFunc("/chat/", chat)
	mux.HandleFunc("/ws/", webSocket)
	mux.HandleFunc("/admin", admin)
	mux.HandleFunc("/csv/", csvDown)
	mux.HandleFunc("/delete/", roomDel)
	mux.HandleFunc("/user", userList)
	mux.HandleFunc("/userDel/", userDel)
}
