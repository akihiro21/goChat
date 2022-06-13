package handler

import (
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

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
