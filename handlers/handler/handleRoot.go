package handler

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func root(w http.ResponseWriter, r *http.Request) {
	tokenCheck(w, r)
	if login := nowLoginBool(w, r); login == false {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if sessionName(w, r) == "admin" {
		http.Redirect(w, r, "/admin", http.StatusFound)
	} else {
		http.Redirect(w, r, "/room", http.StatusFound)
	}
}
