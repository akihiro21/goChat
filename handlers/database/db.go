package database

import (
	"database/sql"
)

type UserOperation interface {
	ReadValue(key string, value string, db *sql.DB) (User, error)
	ReadId(id int, db *sql.DB) (User, error)
	ReadAll(db *sql.DB) []User
	Insert(user *User, db *sql.DB)
	RoomUpdate(id int, userName string, db *sql.DB) error
	Delete(name string, db *sql.DB)
	PassCheck(userName string, userPassword string, db *sql.DB) error
}

type RoomOperation interface {
	ReadValue(key string, value string, db *sql.DB) (Room, error)
	ReadId(id int, db *sql.DB) (Room, error)
	ReadAll(admin bool, db *sql.DB) []Room
	Insert(room Room, db *sql.DB)
	UserUpdate(key string, id int, roomName string, db *sql.DB) error
	Delete(room string, db *sql.DB)
}

type MessageOperation interface {
	ReadAll(room string, db *sql.DB) []Message
	Insert(message *Message, db *sql.DB)
	Delete(room string, db *sql.DB)
}
