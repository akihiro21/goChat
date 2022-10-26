package database

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

func (d *messageDatabase) ReadAll(room string, db *sql.DB) []Message {
	var oneMessage Message
	var messages []Message

	prep, err := db.Prepare("select * from message WHERE room = ?;")
	defer prep.Close()
	if err != nil {
		log.Println(err)
		return nil
	}

	rows, err := prep.Query(room)
	if err != nil {
		log.Println(err)
		return nil
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

func (d *messageDatabase) Insert(message *Message, db *sql.DB) {
	prep, err := db.Prepare("INSERT INTO message(message,room,user) VALUES(?,?,?)")
	defer prep.Close()
	if err != nil {
		fmt.Println(err)
	}

	_, err = prep.Exec(message.Message, message.Room, message.UserName)
	if err != nil {
		fmt.Println(err)
	}
}

func (d *messageDatabase) Delete(room string, db *sql.DB) {
	delete, err := db.Prepare("DELETE FROM message WHERE room = ? ")
	defer delete.Close()
	if err != nil {
		log.Println(err)
	}
	_, err = delete.Exec(room)
	defer delete.Close()
	if err != nil {
		log.Println(err)
	}
}
