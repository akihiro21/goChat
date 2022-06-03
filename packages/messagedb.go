package packages

import (
	"database/sql"
	"log"
)

type Message struct {
	id       int
	Message  string
	Room     string
	UserName string
}

func (message *Message) ReadAll(room string, db *sql.DB) []Message {
	var oneMessage Message
	var messages []Message
	prep, err := db.Prepare("select * from message WHERE room = ?;")
	defer prep.Close()
	if err != nil {
		log.Println(err)
	}

	rows, err := prep.Query(room)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&oneMessage.id, &oneMessage.Message, &oneMessage.Room, &oneMessage.UserName)
		if err != nil {
			log.Println(err)
		}
		messages = append(messages, oneMessage)
	}

	return messages
}

func (message *Message) Insert(db *sql.DB) {
	prep, err := db.Prepare("INSERT INTO message(message,room,user) VALUES(?,?,?)")
	defer prep.Close()
	if err != nil {
		log.Println(err)
	}

	_, err = prep.Exec(message.Message, message.Room, message.UserName)
	if err != nil {
		log.Println(err)
	}
}

func (message *Message) Delete(name string, db *sql.DB) {
	delete, err := db.Prepare("DELETE FROM message WHERE room = ? ")
	defer delete.Close()
	if err != nil {
		log.Println(err)
	}
	_, err = delete.Exec(name)
	defer delete.Close()
	if err != nil {
		log.Println(err)
	}
}
