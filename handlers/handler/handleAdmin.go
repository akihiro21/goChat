package handlers

import (
	"encoding/csv"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func admin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if noSession(w, r) {
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			if sessionName(w, r) == "admin" {
				t := templates["admin"]
				tokenCheck(w, r)
				roomAll := roomDB.readAll(true, db)
				if err := t.Execute(w, struct {
					Css   string
					Js    string
					Alert string
					Login bool
					Room  []Room
					Token string
				}{Css: "room", Js: "room", Alert: msg.Message, Login: noSession(w, r), Room: roomAll, Token: token(w, r)}); err != nil {
					log.Printf("failed to execute template: %v", err)
					msg.Message = ""
				}
			} else {
				http.Redirect(w, r, "/room", http.StatusFound)
			}
		}
	}
}

func userList(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if noSession(w, r) {
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			if sessionName(w, r) == "admin" {
				t := templates["user"]
				tokenCheck(w, r)
				userAll := userDB.readAll(db)
				if err := t.Execute(w, struct {
					Css   string
					Js    string
					Alert string
					Login bool
					User  []User
					Token string
				}{Css: "room", Js: "room", Alert: msg.Message, Login: noSession(w, r), User: userAll[1:], Token: token(w, r)}); err != nil {
					log.Printf("failed to execute template: %v", err)
					msg.Message = ""
				}
			} else {
				http.Redirect(w, r, "/room", http.StatusFound)
			}
		}
	}
}

func userDel(w http.ResponseWriter, r *http.Request) {
	if noSession(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else if r.Method == "GET" {
		tokenCheck(w, r)
		ep := strings.TrimPrefix(r.URL.Path, "/userDel")
		_, name := filepath.Split(ep)
		if name != "" {
			userDB.delete(name, db)
		}
	}
	http.Redirect(w, r, "/user", http.StatusFound)
}

func csvDown(w http.ResponseWriter, r *http.Request) {
	if noSession(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else if r.Method == "GET" {
		tokenCheck(w, r)
		ep := strings.TrimPrefix(r.URL.Path, "/csv")
		_, name := filepath.Split(ep)
		if name != "" {
			chats := MessageDB.readAll(name, db)
			list := make([][]string, len(chats))
			for i := range list {
				list[i] = make([]string, 2)
			}
			for i := 0; i < len(list); i++ {
				list[i][0] = chats[i].UserName
				list[i][1] = chats[i].Message
			}

			c := csv.NewWriter(w)
			c.WriteAll(list)
			c.Flush()

			if err := c.Error(); err != nil {
				log.Fatal(err)
			}

		}
	}
	http.Redirect(w, r, "/admin", http.StatusFound)

}

func roomDel(w http.ResponseWriter, r *http.Request) {
	if noSession(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else if r.Method == "GET" {
		tokenCheck(w, r)
		ep := strings.TrimPrefix(r.URL.Path, "/delete")
		_, name := filepath.Split(ep)
		if name != "" {
			roomDB.delete(name, db)
			MessageDB.delete(name, db)
		}
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}
