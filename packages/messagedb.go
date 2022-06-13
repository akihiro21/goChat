package packages

import (
	"database/sql"
	"fmt"
	"log"
)

type Message struct {
	id       int
	Message  string
	Room     string
	UserName string
}

type messageDatabase struct{}

func NewMessageDB() MessageOperation {
	var messageDb messageDatabase
	db := MessageOperation(&messageDb)
	return db
}

func (d *messageDatabase) readAll(room string, db *sql.DB) []Message {
	var oneMessage Message
	var messages []Message

	prep, err := db.Prepare("select * from message WHERE room = ?;")
	defer prep.Close()
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	rows, err := prep.Query(room)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&oneMessage.id, &oneMessage.Message, &oneMessage.Room, &oneMessage.UserName)
		if err != nil {
			log.Println(err.Error())
		}
		messages = append(messages, oneMessage)
	}

	return messages
}

func (d *messageDatabase) insert(message *Message, db *sql.DB) {
	prep, err := db.Prepare("INSERT INTO message(message,room,user) VALUES(?,?,?)")
	defer prep.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = prep.Exec(message.Message, message.Room, message.UserName)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (d *messageDatabase) delete(room string, db *sql.DB) {
	delete, err := db.Prepare("DELETE FROM message WHERE room = ? ")
	defer delete.Close()
	if err != nil {
		log.Println(err.Error())
	}
	_, err = delete.Exec(room)
	defer delete.Close()
	if err != nil {
		log.Println(err.Error())
	}
}
