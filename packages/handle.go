package packages

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Msg struct {
	Message string
}

var (
	tokens    []string
	chatMess  = Message{}
	msg       = Msg{}
	account   = User{}
	rooms     = Room{}
	db        *sql.DB
	templates = make(map[string]*template.Template)
)

func HandleInit(mux *http.ServeMux) {
	go H.run()
	templates["login"] = loadTemplate("login")
	templates["create"] = loadTemplate("create")
	templates["room"] = loadTemplate("room")
	templates["chat"] = loadTemplate("chat")
	mux.HandleFunc("/", root)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/create", create)
	mux.HandleFunc("/room", room)
	mux.HandleFunc("/chat/", chat)
	mux.HandleFunc("/ws/", webSocket)
}

func DbInit(myDb *sql.DB) {
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
		http.Redirect(w, r, "/login", 302)
	} else {
		http.Redirect(w, r, "/room", 302)
	}

}

func login(w http.ResponseWriter, r *http.Request) {
	t := templates["login"]
	if r.Method == "GET" {
		tokenCheck(w, r)
		if err := t.Execute(w, struct {
			Css   string
			Alert string
			Login bool
			Token string
		}{Css: "login", Alert: msg.Message, Login: noSession(w, r), Token: token(w, r)}); err != nil {
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
				http.Redirect(w, r, "/login", 302)
			} else if len(r.Form.Get("password")) == 0 {
				msg.Message = "Passwordが入力されていません。"
				http.Redirect(w, r, "/login", 302)
			} else {
				account = User{
					name:     r.Form.Get("username"),
					password: r.Form.Get("password"),
				}
				if account.UserCheck(db) == nil {
					if account.PassCheck(db) == nil {
						loginSession(account.name, w, r)
						http.Redirect(w, r, "/room", 302)
					} else {
						msg.Message = "パスワードが違います。"
						http.Redirect(w, r, "/login", 302)
					}
				} else {
					msg.Message = "アカウントは存在しません。"
					http.Redirect(w, r, "/login", 302)
				}
			}
		}
		http.Redirect(w, r, "/login", 302)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	deleteSession(w, r)
	http.Redirect(w, r, "/login", 302)
}

func create(w http.ResponseWriter, r *http.Request) {
	t := templates["create"]
	if r.Method == "GET" {
		tokenCheck(w, r)
		if !noSession(w, r) {
			http.Redirect(w, r, "/room", 302)
		} else {
			if err := t.Execute(w, struct {
				Css   string
				Alert string
				Login bool
				Token string
			}{Css: "login", Alert: msg.Message, Login: noSession(w, r), Token: token(w, r)}); err != nil {
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
				http.Redirect(w, r, "/create", 302)
			} else if len(r.Form.Get("password")) == 0 {
				msg.Message = "Passwordが入力されていません。"
				http.Redirect(w, r, "/create", 302)
			} else {
				account = User{
					name:     r.Form.Get("username"),
					password: r.Form.Get("password"),
					room:     0,
				}
				if account.UserCheck(db) == nil {
					msg.Message = "そのユーザネームは既に存在します。"
					http.Redirect(w, r, "/create", 302)
				} else {
					account.Insert(db)
					loginSession(account.name, w, r)
					http.Redirect(w, r, "/room", 302)
				}
			}
		}
		http.Redirect(w, r, "/create", 302)
	}
}

func room(w http.ResponseWriter, r *http.Request) {
	if noSession(w, r) {
		http.Redirect(w, r, "/login", 302)
	} else {
		if r.Method == "GET" {
			t := templates["room"]
			tokenCheck(w, r)
			roomAll := rooms.ReadAll(db)
			if err := t.Execute(w, struct {
				Css   string
				Alert string
				Login bool
				Room  []Room
				Token string
			}{Css: "room", Alert: msg.Message, Login: noSession(w, r), Room: roomAll, Token: token(w, r)}); err != nil {
				log.Printf("failed to execute template: %v", err)
			}
			msg.Message = ""
		} else if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				log.Println(err)
			}

			t := r.Form.Get("token")
			if t == token(w, r) {
				session, err := sessionStore.Get(r, SessionName)
				if err != nil {
					handleSessionError(w, err)
					return
				}
				account.name, _ = session.Values["username"].(string)
				account.ReadOne(db)

				t := r.Form.Get("date")
				rooms = Room{
					Name:    r.Form.Get("name"),
					Date:    t,
					UserNum: 0,
					user1:   account.id,
					user2:   0,
				}
				if rooms.nameCheck(db) == nil {
					msg.Message = "そのルーム名は既に存在します。"
					http.Redirect(w, r, "/room", 302)
				}
				rooms.Insert(db)
				http.Redirect(w, r, "/room", 302)
			}
		}
	}
}

func chat(w http.ResponseWriter, r *http.Request) {
	if noSession(w, r) {
		http.Redirect(w, r, "/login", 302)
	} else if r.Method == "GET" {
		tokenCheck(w, r)
		ep := strings.TrimPrefix(r.URL.Path, "/chat")
		_, name := filepath.Split(ep)
		if name != "" {
			rooms.Name = name
			if rooms.nameCheck(db) == nil {
				t := templates["chat"]
				msg.Message = ""
				chats := chatMess.ReadAll(name, db)
				if err := t.Execute(w, struct {
					Css    string
					Alert  string
					Login  bool
					Chat   []Message
					Room   string
					MyName string
					Token  string
				}{Css: "chat", Alert: msg.Message, Login: noSession(w, r), Chat: chats, Room: rooms.Name, MyName: sessionName(w, r), Token: token(w, r)}); err != nil {
					log.Printf("failed to execute template: %v", err)
				}
				msg.Message = ""
			} else {
				http.Redirect(w, r, "/room", 302)
			}
		} else {
			http.Redirect(w, r, "/room", 302)
		}
	} else if r.Method == "POST" {
		ep := strings.TrimPrefix(r.URL.Path, "/chat")
		_, name := filepath.Split(ep)
		if name != "" {
			if err := r.ParseForm(); err != nil {
				log.Println(err)
			}
			t := r.Form.Get("token")
			if t == token(w, r) {
				t := r.Form.Get("message")
				mes := Message{
					Message:  t,
					Room:     name,
					UserName: sessionName(w, r),
				}
				mes.Insert(db)
			}
			http.Redirect(w, r, "/chat/"+name, 302)
		}
	}
}

func webSocket(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		ep := strings.TrimPrefix(r.URL.Path, "/chat")
		_, name := filepath.Split(ep)
		if name != "" {
			serveWs(w, r, name, sessionName(w, r))
		} else {
			http.Redirect(w, r, "/chat/"+name, 302)
		}
	}
}
