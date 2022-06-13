package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Msg struct {
	Message string
}

var (
	tokens    []string
	msg       = Msg{}
	db        *sql.DB
	userDB    = NewUserDB()
	roomDB    = NewRoomDB()
	MessageDB = NewMessageDB()
	templates = make(map[string]*template.Template)
)

func Init(myDb *sql.DB) {
	db = myDb
}

func loadTemplate(name string) *template.Template {
	t, err := template.ParseFiles(
		"web/templates/"+name+".html",
		"web/templates/_header.html",
		"web/templates/_footer.html",
	)
	if err != nil {
		log.Fatalf("template error: %v", err)
	}

	return t
}

func root(w http.ResponseWriter, r *http.Request) {
	tokenCheck(w, r)
	if noSession(w, r) {
		http.Redirect(w, r, "/create", http.StatusFound)
	} else {
		if sessionName(w, r) == "admin" {
			http.Redirect(w, r, "/admin", http.StatusFound)
		} else {
			http.Redirect(w, r, "/room", http.StatusFound)
		}
	}

}
