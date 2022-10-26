package handler

import (
	"encoding/csv"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/akihiro21/goChat/handlers/database"
)

func admin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if login := nowLoginBool(w, r); login == false {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if sessionName(w, r) == "admin" {
			t := templates["admin"]
			tokenCheck(w, r)
			roomAll := roomDB.ReadAll(true, db)
			if err := t.Execute(w, struct {
				Css   string
				Js    string
				Alert string
				Room  []database.Room
				Token string
				Login bool
			}{Css: "room", Js: "room", Alert: msg.Message, Room: roomAll, Token: token(w, r), Login: nowLoginBool(w, r)}); err != nil {
				log.Printf("failed to execute template: %v", err)
				msg.Message = ""
			}
		} else {
			http.Redirect(w, r, "/room", http.StatusFound)
		}
	}
}

func userList(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if login := nowLoginBool(w, r); login == false {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if sessionName(w, r) != "admin" {
			http.Redirect(w, r, "/room", http.StatusFound)
			return
		}
		t := templates["user"]
		tokenCheck(w, r)
		userAll := userDB.ReadAll(db)
		if err := t.Execute(w, struct {
			Css   string
			Js    string
			Alert string
			User  []database.User
			Token string
			Login bool
		}{Css: "room", Js: "room", Alert: msg.Message, User: userAll[1:], Token: token(w, r), Login: nowLoginBool(w, r)}); err != nil {
			log.Printf("failed to execute template: %v", err)
			msg.Message = ""
		}
	}
}

func userDel(w http.ResponseWriter, r *http.Request) {
	if login := nowLoginBool(w, r); login == false {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method == "GET" {
		tokenCheck(w, r)
		ep := strings.TrimPrefix(r.URL.Path, "/userDel")
		_, name := filepath.Split(ep)
		if name != "" {
			userDB.Delete(name, db)
		}
	}
	http.Redirect(w, r, "/user", http.StatusFound)
}

func csvDown(w http.ResponseWriter, r *http.Request) {
	if login := nowLoginBool(w, r); login == false {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method == "GET" {
		tokenCheck(w, r)
		ep := strings.TrimPrefix(r.URL.Path, "/csv")
		_, name := filepath.Split(ep)
		if name != "" {
			chats := MessageDB.ReadAll(name, db)
			list := make([][]string, len(chats))
			for i := range list {
				list[i] = make([]string, 2)
			}
			for i := 0; i < len(list); i++ {
				list[i][0] = chats[i].UserName
				list[i][1] = chats[i].Message
			}

			c := csv.NewWriter(w)
			err := c.WriteAll(list)
			if err != nil {
				log.Println(err)
			}
			c.Flush()

			if err := c.Error(); err != nil {
				log.Println(err)
			}

		}
	}
	http.Redirect(w, r, "/admin", http.StatusFound)

}

func roomDel(w http.ResponseWriter, r *http.Request) {
	if login := nowLoginBool(w, r); login == false {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method == "GET" {
		tokenCheck(w, r)
		ep := strings.TrimPrefix(r.URL.Path, "/delete")
		_, name := filepath.Split(ep)
		if name != "" {
			roomDB.Delete(name, db)
			MessageDB.Delete(name, db)
		}
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}
