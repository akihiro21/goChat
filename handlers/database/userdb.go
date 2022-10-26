package database

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int
	Name     string
	Password string
	Room     int
}

type userDatabase struct{}

func NewUserDB() UserOperation {
	var userDB userDatabase
	db := UserOperation(&userDB)
	return db
}

//全てのユーザデータを取得
func (d *userDatabase) ReadAll(db *sql.DB) []User {
	var user User
	var users []User
	rows, err := db.Query("select * from user;")
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Password, &user.Room)
		if err != nil {
			log.Println(err)
		}
		users = append(users, user)
	}
	rows.Close()
	return users
}

//特定のユーザデータを取得
func (d *userDatabase) ReadValue(key string, value string, db *sql.DB) (User, error) {
	var (
		sql  string
		user User
	)

	switch key {
	case "name":
		sql = "SELECT * FROM user WHERE name = ? LIMIT 1"
	}

	prep, err := db.Prepare(sql)
	defer prep.Close()
	if err != nil {
		log.Println(err)
	}

	err = prep.QueryRow(value).Scan(&user.Id, &user.Name, &user.Password, &user.Room)
	if err != nil {
		log.Println(err)
	}

	return user, err
}

func (d *userDatabase) ReadId(id int, db *sql.DB) (User, error) {
	var user User
	prep, err := db.Prepare("SELECT * FROM user WHERE id = ? LIMIT 1")
	defer prep.Close()
	if err != nil {
		log.Println(err)
	}

	err = prep.QueryRow(id).Scan(&user.Id, &user.Name, &user.Password, &user.Room)
	if err != nil {
		log.Println(err)
	}

	return user, err
}

//新しいデータの追加
func (d *userDatabase) Insert(user *User, db *sql.DB) {
	ins, err := db.Prepare("INSERT INTO user(name,password,room) VALUES(?,?,?)")
	defer ins.Close()
	if err != nil {
		log.Println(err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		log.Println(err)
	}

	_, err = ins.Exec(user.Name, string(hash), user.Room)
	if err != nil {
		log.Println(err)
	}
}

//パスワードの変更
func (d *userDatabase) RoomUpdate(id int, userName string, db *sql.DB) error {
	var (
		room Room
	)
	roomDB := NewRoomDB()

	_, err := d.ReadValue("name", userName, db)
	if err != nil {
		log.Println(err)
		return err
	}

	room, err = roomDB.ReadId(id, db)
	if err != nil {
		log.Println(err)
		return err
	}
	upd, err := db.Prepare("UPDATE user SET room = ? WHERE ( name = ? ) LIMIT 1")
	defer upd.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = upd.Exec(room.Id, userName)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (d *userDatabase) PassCheck(userName string, userPassword string, db *sql.DB) error {
	var user User
	exist, err := db.Prepare("SELECT * FROM user WHERE name = ? LIMIT 1")
	defer exist.Close()
	if err != nil {
		log.Println("dbOpen", err)
	}
	err = exist.QueryRow(userName).Scan(&user.Id, &user.Name, &user.Password, &user.Room)
	if err == nil {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPassword))
		log.Println("PassCheck", err)
		return err
	}
	return err
}

//データの消去
func (d *userDatabase) Delete(name string, db *sql.DB) {
	delete, err := db.Prepare("DELETE FROM user WHERE name = ? ")
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
