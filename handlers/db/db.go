package handlers

import (
	"database/sql"
)

type UserOperation interface {
	readValue(key string, value string, db *sql.DB) (User, error)
	readId(id int, db *sql.DB) (User, error)
	readAll(db *sql.DB) []User
	insert(user *User, db *sql.DB)
	roomUpdate(id int, userName string, db *sql.DB) error
	delete(name string, db *sql.DB)
	passCheck(userName string, userPassword string, db *sql.DB) error
}

type RoomOperation interface {
	readValue(key string, value string, db *sql.DB) (Room, error)
	readId(id int, db *sql.DB) (Room, error)
	readAll(admin bool, db *sql.DB) []Room
	insert(room Room, db *sql.DB)
	userUpdate(key string, id int, roomName string, db *sql.DB) error
	delete(room string, db *sql.DB)
}

type MessageOperation interface {
	readAll(room string, db *sql.DB) []Message
	insert(message *Message, db *sql.DB)
	delete(room string, db *sql.DB)
}
