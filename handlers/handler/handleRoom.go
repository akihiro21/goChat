package handlers

import (
	"log"
	"net/http"
)

func room(w http.ResponseWriter, r *http.Request) {
	if noSession(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		if r.Method == "GET" {
			tokenCheck(w, r)
			name := sessionName(w, r)
			account, _ := userDB.readValue("name", name, db)
			if account.room != 0 {
				room, err := roomDB.readId(account.room, db)
				if err != nil {
					log.Println("roomDB err" + err.Error())
				}
				url := "/chat/" + room.Name
				http.Redirect(w, r, url, http.StatusFound)
			}
			t := templates["room"]
			rooms := roomDB.readAll(false, db)
			if err := t.Execute(w, struct {
				Css   string
				Js    string
				Alert string
				Login bool
				Room  []Room
				Token string
			}{Css: "room", Js: "room", Alert: msg.Message, Login: noSession(w, r), Room: rooms, Token: token(w, r)}); err != nil {
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
				account, _ := userDB.readValue("name", name, db)

				t := r.Form.Get("date")
				room := Room{
					Name:    r.Form.Get("name"),
					Date:    t,
					UserNum: 0,
					user1:   account.id,
					user2:   0,
				}
				if _, err := roomDB.readValue("name", room.Name, db); err == nil {
					msg.Message = "そのルーム名は既に存在します。"
					http.Redirect(w, r, "/room", http.StatusFound)
				}
				roomDB.insert(room, db)
				http.Redirect(w, r, "/room", http.StatusFound)
			}
		}
	}
}
