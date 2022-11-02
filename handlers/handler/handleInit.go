package handler

import (
	"database/sql"
	"encoding/csv"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/akihiro21/goChat/handlers/database"
)

type Msg struct {
	Message string
}

type Scenario struct {
	JikDay1 []string
	YwkDay1 []string
	JikDay2 []string
	YwkDay2 []string
	JikDay3 []string
	YwkDay3 []string
}

var (
	tokens    []string
	msg       = Msg{}
	db        *sql.DB
	userDB    = database.NewUserDB()
	roomDB    = database.NewRoomDB()
	MessageDB = database.NewMessageDB()
	templates = make(map[string]*template.Template)
	scenario1 = loadCsv("scenario1")
	scenario2 = loadCsv("scenario2")
)

func Init(mux *http.ServeMux, myDb *sql.DB) {
	go H.run()
	db = myDb
	SessionInit()
	templatesInit()
	HandlerInit(mux)
}

func templatesInit() {
	templates["login"] = loadTemplate("login")
	templates["create"] = loadTemplate("create")
	templates["room"] = loadTemplate("room")
	templates["chat"] = loadTemplate("chat")
	templates["adminChat"] = loadTemplate("adminChat")
	templates["admin"] = loadTemplate("admin")
	templates["user"] = loadTemplate("user")
}

func HandlerInit(mux *http.ServeMux) {
	mux.HandleFunc("/", root)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/create", create)
	mux.HandleFunc("/room", room)
	mux.HandleFunc("/chat/", chat)
	mux.HandleFunc("/ws/", webSocket)
	mux.HandleFunc("/admin", admin)
	mux.HandleFunc("/csv/", csvDown)
	mux.HandleFunc("/delete/", roomDel)
	mux.HandleFunc("/user", userList)
	mux.HandleFunc("/userDel/", userDel)
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

func loadCsv(name string) Scenario {
	file, err := os.Open("web/csv/" + name + ".csv")
	if err != nil {
		log.Printf("csv load error: %v", err)
	}
	defer file.Close()
	csvReader := csv.NewReader(file)

	scenario := Scenario{}
	for {
		line, err := csvReader.Read()
		if err != nil {
			log.Println(err)
			break
		}

		scenario.JikDay1 = append(scenario.JikDay1, line[0])
		scenario.YwkDay1 = append(scenario.YwkDay1, line[1])
		scenario.JikDay2 = append(scenario.JikDay2, line[2])
		scenario.YwkDay2 = append(scenario.YwkDay2, line[3])
		scenario.JikDay3 = append(scenario.JikDay3, line[4])
		scenario.YwkDay3 = append(scenario.YwkDay3, line[5])
	}

	return scenario
}
