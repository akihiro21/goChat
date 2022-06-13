package handlers

import "net/http"

func HandleInit(mux *http.ServeMux) {
	go H.run()
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
