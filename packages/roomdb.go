package packages

import (
	"database/sql"
	"log"
	"time"
)

type Room struct {
	Id      int
	Name    string
	Date    string
	UserNum int
	user1   int
	user2   int
}

//全てのユーザデータを取得
func (room *Room) ReadAll(db *sql.DB) []Room {
	var (
		oneRoom Room
		rooms   []Room
		date    time.Time
	)
	rows, err := db.Query("select * from room order by date;")
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		err = rows.Scan(&oneRoom.Id, &oneRoom.Name, &date, &oneRoom.user1, &oneRoom.user2)
		if err != nil {
			log.Println(err)
		}
		oneRoom.Date = date.Format("2006-01-02 15:04")
		oneRoom.UserNum = oneRoom.num()
		rooms = append(rooms, oneRoom)
	}
	rows.Close()
	return rooms
}

func (room *Room) Insert(db *sql.DB) {
	ins, err := db.Prepare("INSERT INTO room(name,date,userId1,userId2) VALUES(?,?,?,?)")
	defer ins.Close()
	if err != nil {
		log.Println(err)
	}

	_, err = ins.Exec(room.Name, room.Date, room.user1, room.user2)
	if err != nil {
		log.Println(err)
	}
}

func (room *Room) nameCheck(db *sql.DB) error {
	exist, err := db.Prepare("SELECT * FROM room WHERE name = ? LIMIT 1")
	defer exist.Close()
	if err != nil {
		panic(err)
	}
	err = exist.QueryRow(room.Name).Scan(&room.Id, &room.Name, &room.Date, &room.user1, &room.user2)
	return err
}

func (room *Room) idCheck(db *sql.DB) error {
	exist, err := db.Prepare("SELECT * FROM room WHERE id = ? LIMIT 1")
	defer exist.Close()
	if err != nil {
		panic(err)
	}
	err = exist.QueryRow(room.Id).Scan(&room.Id, &room.Name, &room.Date, &room.user1, &room.user2)
	return err
}

func (room *Room) UpdUser1(db *sql.DB) {
	upd, err := db.Prepare("UPDATE room SET userId1 = ? WHERE ( name = ? ) LIMIT 1")
	defer upd.Close()
	if err != nil {
		log.Println(err)
	}
	_, err = upd.Exec(room.user1, room.Name)
	if err != nil {
		log.Println(err)
	}
}

func (room *Room) UpdUser2(db *sql.DB) {
	upd, err := db.Prepare("UPDATE room SET userId2 = ? WHERE ( name = ? ) LIMIT 1")
	defer upd.Close()
	if err != nil {
		log.Println(err)
	}
	_, err = upd.Exec(room.user2, room.Name)
	if err != nil {
		log.Println(err)
	}
}

func (room *Room) Delete(db *sql.DB) {
	delete, err := db.Prepare("DELETE FROM room WHERE name = ? ")
	defer delete.Close()
	if err != nil {
		log.Println(err)
	}
	_, err = delete.Exec(room.Name)
	defer delete.Close()
	if err != nil {
		log.Println(err)
	}
}

func (room *Room) num() int {
	if room.user2 != 0 && room.user1 != 0 {
		return 2
	} else if room.user2 == 0 && room.user1 == 0 {
		return 0
	} else {
		return 1
	}
}
