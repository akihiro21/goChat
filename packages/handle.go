package packages

import (
	"database/sql"
	"encoding/csv"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

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

func HandleInit(mux *http.ServeMux) {
	go H.run()
	templates["login"] = loadTemplate("login")
	templates["create"] = loadTemplate("create")
	templates["room"] = loadTemplate("room")
	templates["chat"] = loadTemplate("chat")
	templates["adminChat"] = loadTemplate("adminChat")
	templates["admin"] = loadTemplate("admin")
	templates["user"] = loadTemplate("user")
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

func chat(w http.ResponseWriter, r *http.Request) {
	if noSession(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else if r.Method == "GET" {
		tokenCheck(w, r)
		ep := strings.TrimPrefix(r.URL.Path, "/chat")
		_, name := filepath.Split(ep)
		if name != "" {
			account, _ := userDB.readValue("name", name, db)
			if room, err := roomDB.readValue("name", name, db); err == nil {
				if account.Name != "admin" {
					if err := roomDB.userUpdate("userId1", account.id, room.Name, db); err != nil {
						if err := roomDB.userUpdate("userId2", account.id, room.Name, db); err != nil {
							msg.Message = "この部屋は満員です。"
							http.Redirect(w, r, "/room", http.StatusFound)
						}
					}
					if err := userDB.roomUpdate(room.Id, account.Name, db); err != nil {
						log.Println(err.Error())
					}
				}

				if sessionName(w, r) == "admin" {
					t := templates["adminChat"]
					msg.Message = ""
					chats := MessageDB.readAll(name, db)
					if err := t.Execute(w, struct {
						Css    string
						Js     string
						Alert  string
						Login  bool
						Chat   []Message
						Room   string
						MyName string
						User   string
						Token  string
					}{Css: "adminChat", Js: "chat", Alert: msg.Message, Login: noSession(w, r), Chat: chats, Room: room.Name, MyName: sessionName(w, r), User: account.Name, Token: token(w, r)}); err != nil {
						log.Printf("failed to execute template: %v", err)
					}
					msg.Message = ""
				} else {
					t := templates["chat"]
					msg.Message = ""
					chats := MessageDB.readAll(name, db)
					if err := t.Execute(w, struct {
						Css    string
						Js     string
						Alert  string
						Login  bool
						Chat   []Message
						Room   string
						MyName string
						User   string
						Token  string
					}{Css: "chat", Js: "chat", Alert: msg.Message, Login: noSession(w, r), Chat: chats, Room: room.Name, MyName: sessionName(w, r), User: "OrangeBot", Token: token(w, r)}); err != nil {
						log.Printf("failed to execute template: %v", err)
					}
					msg.Message = ""
				}
			} else {
				http.Redirect(w, r, "/room", http.StatusFound)
			}
		} else {
			http.Redirect(w, r, "/room", http.StatusFound)
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
				t := strings.TrimRight(r.Form.Get("message"), "\n")
				if t != "" {
					mes := Message{
						Message:  t,
						Room:     name,
						UserName: sessionName(w, r),
					}
					MessageDB.insert(&mes, db)
				}
			}
			http.Redirect(w, r, "/chat/"+name, http.StatusFound)
		}
	}
}

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

func webSocket(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		ep := strings.TrimPrefix(r.URL.Path, "/chat")
		_, name := filepath.Split(ep)
		if name != "" {
			serveWs(w, r, name, sessionName(w, r))
		} else {
			http.Redirect(w, r, "/chat/"+name, http.StatusFound)
		}
	}
}
