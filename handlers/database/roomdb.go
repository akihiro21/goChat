package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Room struct {
	Id      int
	Name    string
	Date    string
	UserNum int
	User1   int
	User2   int
}

type roomDatabase struct{}

func NewRoomDB() RoomOperation {
	var roomDB roomDatabase
	db := RoomOperation(&roomDB)
	return db
}

//全てのユーザデータを取得
func (d *roomDatabase) ReadAll(admin bool, db *sql.DB) []Room {
	var (
		oneRoom Room
		rooms   []Room
		date    time.Time
		sql     string
	)

	if admin {
		sql = "select * from room order by date;"
	} else {
		sql = "select * from room WHERE userId1 = 0 or userId2 = 0 order by date;"
	}

	rows, err := db.Query(sql)
	if err != nil {
		log.Println("roomReadAll" + err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&oneRoom.Id, &oneRoom.Name, &date, &oneRoom.User1, &oneRoom.User2)
		if err != nil {
			log.Println("roomReadAll" + err.Error())
		}
		oneRoom.Date = date.Format("2006-01-02 15:04")
		oneRoom.UserNum = d.num(&oneRoom)
		rooms = append(rooms, oneRoom)
	}
	rows.Close()
	return rooms
}

func (d *roomDatabase) Insert(room Room, db *sql.DB) {
	ins, err := db.Prepare("INSERT INTO room(name,date,userId1,userId2) VALUES(?,?,?,?)")
	defer ins.Close()
	if err != nil {
		log.Println(err)
	}

	_, err = ins.Exec(&room.Name, &room.Date, &room.User1, &room.User2)
	if err != nil {
		log.Println("Insert room" + err.Error())
	}
}

func (d *roomDatabase) ReadValue(key string, value string, db *sql.DB) (Room, error) {
	var (
		room Room
		sql  string
	)

	switch key {
	case "name":
		sql = "SELECT * FROM room WHERE name = ? LIMIT 1"
	}

	exist, err := db.Prepare(sql)
	defer exist.Close()
	if err != nil {
		log.Println(err)
	}

	err = exist.QueryRow(value).Scan(&room.Id, &room.Name, &room.Date, &room.User1, &room.User2)
	if err != nil {
		log.Println(err)
	}

	return room, err

}

func (d *roomDatabase) ReadId(id int, db *sql.DB) (Room, error) {
	var room Room

	exist, err := db.Prepare("SELECT * FROM room WHERE id = ? LIMIT 1")
	defer exist.Close()
	if err != nil {
		log.Println(err)
	}

	err = exist.QueryRow(id).Scan(&room.Id, &room.Name, &room.Date, &room.User1, &room.User2)
	if err != nil {
		log.Println(err)
	}
	return room, err
}

func (d *roomDatabase) UserUpdate(key string, id int, roomName string, db *sql.DB) error {
	var (
		sql string
	)

	roomOne, err := d.ReadValue("name", roomName, db)
	if err != nil {
		return err
	}

	switch key {
	case "userId1":
		if roomOne.User1 == id {
			return nil
		} else if roomOne.User1 != 0 {
			return fmt.Errorf("Error: %s", "room is crowded")
		}
		sql = "UPDATE room SET userId1 = ? WHERE ( name = ? ) LIMIT 1"

	case "userId2":
		if roomOne.User2 == id {
			return nil
		} else if roomOne.User2 != 0 {
			return fmt.Errorf("Error: %s", "room is crowded")
		}
		sql = "UPDATE room SET userId2 = ? WHERE ( name = ? ) LIMIT 1"
	}

	upd, err := db.Prepare(sql)
	defer upd.Close()
	if err != nil {
		log.Println(err)
	}
	_, err = upd.Exec(id, roomName)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (d *roomDatabase) Delete(room string, db *sql.DB) {
	delete, err := db.Prepare("DELETE FROM room WHERE name = ? ")
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

func (d roomDatabase) num(room *Room) int {
	if room.User2 != 0 && room.User1 != 0 {
		return 2
	} else if room.User2 == 0 && room.User1 == 0 {
		return 0
	} else {
		return 1
	}
}
