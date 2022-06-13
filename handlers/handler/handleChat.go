package handlers

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func chat(w http.ResponseWriter, r *http.Request) {
	if noSession(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else if r.Method == "GET" {
		tokenCheck(w, r)
		ep := strings.TrimPrefix(r.URL.Path, "/chat")
		_, name := filepath.Split(ep)
		if name != "" {
			account, _ := userDB.readValue("name", name, db)
			if room, err := roomDB.readValue("name", name, db); err == nil {
				if account.Name != "admin" {
					if err := roomDB.userUpdate("userId1", account.id, room.Name, db); err != nil {
						if err := roomDB.userUpdate("userId2", account.id, room.Name, db); err != nil {
							msg.Message = "この部屋は満員です。"
							http.Redirect(w, r, "/room", http.StatusFound)
						}
					}
					if err := userDB.roomUpdate(room.Id, account.Name, db); err != nil {
						log.Println(err.Error())
					}
				}

				if sessionName(w, r) == "admin" {
					t := templates["adminChat"]
					msg.Message = ""
					chats := MessageDB.readAll(name, db)
					if err := t.Execute(w, struct {
						Css    string
						Js     string
						Alert  string
						Login  bool
						Chat   []Message
						Room   string
						MyName string
						User   string
						Token  string
					}{Css: "adminChat", Js: "chat", Alert: msg.Message, Login: noSession(w, r), Chat: chats, Room: room.Name, MyName: sessionName(w, r), User: account.Name, Token: token(w, r)}); err != nil {
						log.Printf("failed to execute template: %v", err)
					}
					msg.Message = ""
				} else {
					t := templates["chat"]
					msg.Message = ""
					chats := MessageDB.readAll(name, db)
					if err := t.Execute(w, struct {
						Css    string
						Js     string
						Alert  string
						Login  bool
						Chat   []Message
						Room   string
						MyName string
						User   string
						Token  string
					}{Css: "chat", Js: "chat", Alert: msg.Message, Login: noSession(w, r), Chat: chats, Room: room.Name, MyName: sessionName(w, r), User: "OrangeBot", Token: token(w, r)}); err != nil {
						log.Printf("failed to execute template: %v", err)
					}
					msg.Message = ""
				}
			} else {
				http.Redirect(w, r, "/room", http.StatusFound)
			}
		} else {
			http.Redirect(w, r, "/room", http.StatusFound)
		}
	} else if r.Method == "POST" {
		ep := strings.TrimPrefix(r.URL.Path, "/chat")
		_, name := filepath.Split(ep)
		if name != "" {
			if err := r.ParseForm(); err != nil {
				log.Println(err)
			}
			t := r.Form.Get("token")
			if t == token(w, r) {
				t := strings.TrimRight(r.Form.Get("message"), "\n")
				if t != "" {
					mes := Message{
						Message:  t,
						Room:     name,
						UserName: sessionName(w, r),
					}
					MessageDB.insert(&mes, db)
				}
			}
			http.Redirect(w, r, "/chat/"+name, http.StatusFound)
		}
	}
}

func webSocket(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		ep := strings.TrimPrefix(r.URL.Path, "/chat")
		_, name := filepath.Split(ep)
		if name != "" {
			serveWs(w, r, name, sessionName(w, r))
		} else {
			http.Redirect(w, r, "/chat/"+name, http.StatusFound)
		}
	}
}
