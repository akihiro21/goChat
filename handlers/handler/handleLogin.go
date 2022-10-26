package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/akihiro21/goChat/handlers/database"
)

func login(w http.ResponseWriter, r *http.Request) {
	t := templates["login"]
	if r.Method == "GET" {
		tokenCheck(w, r)
		if login := nowLoginBool(w, r); login == true {
			http.Redirect(w, r, "/room", http.StatusFound)
			return
		}
		if err := t.Execute(w, struct {
			Css   string
			Js    string
			Alert string
			Token string
			Login bool
		}{Css: "login", Js: "", Alert: msg.Message, Token: token(w, r), Login: nowLoginBool(w, r)}); err != nil {
			log.Printf("failed to execute template: %v", err)
		}
		msg.Message = ""

	} else if r.Method == "POST" {
		msg.Message = ""
		if err := r.ParseForm(); err != nil {
			log.Println(err)
		}
		t := r.Form.Get("token")
		if t == token(w, r) {
			if len(r.Form.Get("username")) == 0 {
				msg.Message = "Usernameが入力されていません。"
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			if len(r.Form.Get("password")) == 0 {
				msg.Message = "Passwordが入力されていません。"
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			account := database.User{
				Name:     r.Form.Get("username"),
				Password: r.Form.Get("password"),
			}
			if _, err := userDB.ReadValue("name", account.Name, db); err != nil {
				msg.Message = "アカウントは存在しません。"
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			if userDB.PassCheck(account.Name, account.Password, db) != nil {
				msg.Message = "パスワードが違います。"
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			loginSession(account.Name, w, r)
			if sessionName(w, r) == "admin" {
				http.Redirect(w, r, "/admin", http.StatusFound)
				return
			}
			if sessionName(w, r) != "" {
				http.Redirect(w, r, "/room", http.StatusFound)
			}
		}
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
		if login := nowLoginBool(w, r); login == true {
			http.Redirect(w, r, "/room", http.StatusFound)
			return
		}
		if err := t.Execute(w, struct {
			Css   string
			Js    string
			Alert string
			Token string
			Login bool
		}{Css: "login", Js: "", Alert: msg.Message, Token: token(w, r), Login: nowLoginBool(w, r)}); err != nil {
			log.Printf("failed to execute template: %v", err)
		}
		msg.Message = ""

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
				return
			}
			if len(r.Form.Get("password")) == 0 {
				msg.Message = "Passwordが入力されていません。"
				http.Redirect(w, r, "/create", http.StatusFound)
				return
			}
			account := database.User{
				Name:     r.Form.Get("username"),
				Password: r.Form.Get("password"),
				Room:     0,
			}
			if _, err := userDB.ReadValue("name", account.Name, db); err == nil {
				msg.Message = "そのユーザネームは既に存在します。"
				http.Redirect(w, r, "/create", http.StatusFound)
				return
			}
			userDB.Insert(&account, db)
			data := time.Now().Format("2006-01-02 15:04")
			room := database.Room{
				Name:    account.Name,
				Date:    data,
				UserNum: 0,
				User1:   account.Id,
				User2:   0,
			}
			if _, err := roomDB.ReadValue("name", room.Name, db); err == nil {
				msg.Message = "そのルーム名は既に存在します。"
				http.Redirect(w, r, "/room", http.StatusFound)
				return
			}
			roomDB.Insert(room, db)
			loginSession(account.Name, w, r)
			http.Redirect(w, r, "/chat/"+account.Name, http.StatusFound)

		}
	}
}
