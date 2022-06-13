package handlers

import (
	"log"
	"net/http"
	"time"
)

func login(w http.ResponseWriter, r *http.Request) {
	t := templates["login"]
	if r.Method == "GET" {
		tokenCheck(w, r)
		if err := t.Execute(w, struct {
			Css   string
			Js    string
			Alert string
			Login bool
			Token string
		}{Css: "login", Js: "", Alert: msg.Message, Login: noSession(w, r), Token: token(w, r)}); err != nil {
			log.Printf("failed to execute template: %v", err)
		}
		msg.Message = ""

	} else if r.Method == "POST" {
		msg.Message = ""
		if err := r.ParseForm(); err != nil {
			log.Panicln(err)
		}
		t := r.Form.Get("token")
		if t == token(w, r) {
			if len(r.Form.Get("username")) == 0 {
				msg.Message = "Usernameが入力されていません。"
				http.Redirect(w, r, "/login", http.StatusFound)
			} else if len(r.Form.Get("password")) == 0 {
				msg.Message = "Passwordが入力されていません。"
				http.Redirect(w, r, "/login", http.StatusFound)
			} else {
				account := User{
					Name:     r.Form.Get("username"),
					password: r.Form.Get("password"),
				}

				if _, err := userDB.readValue("name", account.Name, db); err == nil {
					if userDB.passCheck(account.Name, account.password, db) == nil {
						loginSession(account.Name, w, r)
						if sessionName(w, r) == "admin" {
							http.Redirect(w, r, "/admin", http.StatusFound)
						} else {
							http.Redirect(w, r, "/room", http.StatusFound)
						}
					} else {
						msg.Message = "パスワードが違います。"
						http.Redirect(w, r, "/login", http.StatusFound)
					}
				} else {
					msg.Message = "アカウントは存在しません。"
					http.Redirect(w, r, "/login", http.StatusFound)
				}
			}
		}
		http.Redirect(w, r, "/login", 302)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	deleteSession(w, r)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func create(w http.ResponseWriter, r *http.Request) {
	t := templates["create"]
	if r.Method == "GET" {
		tokenCheck(w, r)
		if !noSession(w, r) {
			http.Redirect(w, r, "/room", http.StatusFound)
		} else {
			if err := t.Execute(w, struct {
				Css   string
				Js    string
				Alert string
				Login bool
				Token string
			}{Css: "login", Js: "", Alert: msg.Message, Login: noSession(w, r), Token: token(w, r)}); err != nil {
				log.Printf("failed to execute template: %v", err)
			}
			msg.Message = ""
		}
	} else if r.Method == "POST" {
		msg.Message = ""
		if err := r.ParseForm(); err != nil {
			log.Println(err)
		}
		t := r.Form.Get("token")
		if t == token(w, r) {
			if len(r.Form.Get("username")) == 0 {
				msg.Message = "Usernameが入力されていません。"
				http.Redirect(w, r, "/create", http.StatusFound)
			} else if len(r.Form.Get("password")) == 0 {
				msg.Message = "Passwordが入力されていません。"
				http.Redirect(w, r, "/create", http.StatusFound)
			} else {
				account := User{
					Name:     r.Form.Get("username"),
					password: r.Form.Get("password"),
					room:     0,
				}
				if _, err := userDB.readValue("name", account.Name, db); err == nil {
					msg.Message = "そのユーザネームは既に存在します。"
					http.Redirect(w, r, "/create", http.StatusFound)
				} else {
					userDB.insert(&account, db)
					data := time.Now().Format("2006-01-02 15:04")
					room := Room{
						Name:    account.Name,
						Date:    data,
						UserNum: 0,
						user1:   account.id,
						user2:   0,
					}
					if _, err := roomDB.readValue("name", room.Name, db); err == nil {
						msg.Message = "そのルーム名は既に存在します。"
						http.Redirect(w, r, "/room", http.StatusFound)
					}
					roomDB.insert(room, db)
					loginSession(account.Name, w, r)
					http.Redirect(w, r, "/chat/"+account.Name, http.StatusFound)
				}
			}
		}
		http.Redirect(w, r, "/create", http.StatusFound)
	}
}
