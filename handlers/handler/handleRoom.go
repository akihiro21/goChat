package handler

import (
	"log"
	"net/http"

	"github.com/akihiro21/goChat/handlers/database"
)

func room(w http.ResponseWriter, r *http.Request) {
	if login := nowLoginBool(w, r); login == false {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method == "GET" {
		tokenCheck(w, r)
		name := sessionName(w, r)
		account, _ := userDB.ReadValue("name", name, db)
		if account.Room != 0 {
			room, err := roomDB.ReadId(account.Room, db)
			if err != nil {
				log.Println("roomDB err" + err.Error())
			}
			url := "/chat/" + room.Name
			http.Redirect(w, r, url, http.StatusFound)
		}
		t := templates["room"]
		rooms := roomDB.ReadAll(false, db)
		if err := t.Execute(w, struct {
			Css   string
			Js    string
			Alert string
			Room  []database.Room
			Token string
			Login bool
		}{Css: "room", Js: "room", Alert: msg.Message, Room: rooms, Token: token(w, r), Login: nowLoginBool(w, r)}); err != nil {
			log.Printf("failed to execute template: %v", err)
		}
		msg.Message = ""
	} else if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Println(err)
		}

		t := r.Form.Get("token")
		if t == token(w, r) {
			name := sessionName(w, r)
			account, _ := userDB.ReadValue("name", name, db)

			t := r.Form.Get("date")
			room := database.Room{
				Name:    r.Form.Get("name"),
				Date:    t,
				UserNum: 0,
				User1:   account.Id,
				User2:   0,
			}
			if _, err := roomDB.ReadValue("name", room.Name, db); err == nil {
				msg.Message = "そのルーム名は既に存在します。"
				http.Redirect(w, r, "/room", http.StatusFound)
			}
			roomDB.Insert(room, db)
			http.Redirect(w, r, "/room", http.StatusFound)
		}
	}
}
