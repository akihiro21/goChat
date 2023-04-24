package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/akihiro21/goChat/handlers/database"
)

type SwitchBot struct {
	Command     string `json:"command"`
	Parameter   string `json:"parameter"`
	CommandType string `json:"commandType"`
}

var (
	ACCESS_TOKEN = ""
	API_BASE_URL = "https://api.switch-bot.com"
)

func chat(w http.ResponseWriter, r *http.Request) {
	if login := nowLoginBool(w, r); login == false {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method == "GET" {
		tokenCheck(w, r)
		ep := strings.TrimPrefix(r.URL.Path, "/chat")
		_, name := filepath.Split(ep)
		if name != "" {
			account, _ := userDB.ReadValue("name", name, db)
			room, err := roomDB.ReadValue("name", name, db)
			if err != nil {
				http.Redirect(w, r, "/room", http.StatusFound)
			}
			if account.Name != "admin" {
				if err := roomDB.UserUpdate("userId1", account.Id, room.Name, db); err != nil {
					if err := roomDB.UserUpdate("userId2", account.Id, room.Name, db); err != nil {
						msg.Message = "この部屋は満員です。"
						http.Redirect(w, r, "/room", http.StatusFound)
					}
				}
				if err := userDB.RoomUpdate(room.Id, account.Name, db); err != nil {
					log.Println(err)
				}
			}

			if sessionName(w, r) == "admin" {
				t := templates["adminChat"]
				msg.Message = ""
				chats := MessageDB.ReadAll(name, db)
				if err := t.Execute(w, struct {
					Css       string
					Js        string
					Alert     string
					Chat      []database.Message
					Room      string
					MyName    string
					User      string
					Token     string
					Scenario1 Scenario
					Scenario2 Scenario
					Login     bool
				}{Css: "adminChat", Js: "chat", Alert: msg.Message, Chat: chats, Room: room.Name, MyName: sessionName(w, r), User: account.Name, Token: token(w, r), Scenario1: csvName(scenario1, account.Name), Scenario2: csvName(scenario2, account.Name), Login: nowLoginBool(w, r)}); err != nil {
					log.Printf("failed to execute template: %v", err)
				}
				msg.Message = ""
			} else {
				t := templates["chat"]
				msg.Message = ""
				chats := MessageDB.ReadAll(name, db)
				if err := t.Execute(w, struct {
					Css    string
					Js     string
					Alert  string
					Chat   []database.Message
					Room   string
					MyName string
					User   string
					Token  string
					Login  bool
				}{Css: "chat", Js: "chat", Alert: msg.Message, Chat: chats, Room: room.Name, MyName: sessionName(w, r), User: "Orange", Token: token(w, r), Login: nowLoginBool(w, r)}); err != nil {
					log.Printf("failed to execute template: %v", err)
				}
				msg.Message = ""
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
					mes := database.Message{
						Message:  t,
						Room:     name,
						UserName: sessionName(w, r),
					}
					MessageDB.Insert(&mes, db)

					if sessionName(w, r) == "admin" {
						var command string
						if strings.Contains(t, "つけます") {
							command = "turnOn"
						} else if strings.Contains(t, "けします") {
							command = "turnOff"
						}
						if strings.Contains(t, "テレビ") {
							sbJson(command, "")
						} else if strings.Contains(t, "扇風機") {
							sbJson("電源", "")
						} else if strings.Contains(t, "ライト") {
							sbJson("電源", "")
						}
					}

				}
			}
		}
		http.Redirect(w, r, "/chat/"+name, http.StatusFound)
	}
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

func sbJson(command string, deviceID string) {
	var sbData SwitchBot
	if command == "電源" {
		sbData.Command = command
		sbData.Parameter = "default"
		sbData.CommandType = "customize"
	} else {
		sbData.Command = command
		sbData.Parameter = "default"
		sbData.CommandType = "command"
	}
	jsonData, _ := json.Marshal(sbData)
	err := HttpPost(API_BASE_URL+"/v1.0/devices/"+deviceID+"/commands", jsonData)
	if err != nil {
		log.Println(err)
	}
}

func HttpPost(url string, json []byte) error {
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(json),
	)
	if err != nil {
		return err
	}

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", ACCESS_TOKEN)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray))

	return err
}

func csvName(scenario Scenario, name string) Scenario {
	scenario.JikDay1 = replaceName(scenario.JikDay1, name)
	scenario.YwkDay1 = replaceName(scenario.YwkDay1, name)
	scenario.JikDay2 = replaceName(scenario.JikDay2, name)
	scenario.YwkDay2 = replaceName(scenario.YwkDay2, name)
	scenario.JikDay3 = replaceName(scenario.JikDay3, name)
	scenario.YwkDay3 = replaceName(scenario.YwkDay3, name)

	return scenario
}

func replaceName(scenario []string, name string) []string {
	var ans []string
	for _, v := range scenario {
		ans = append(ans, strings.Replace(v, "UserName", name, -1))
	}

	return ans
}
