package packages

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id       int
	Name     string
	password string
	room     int
}

//全てのユーザデータを取得
func (user *User) readAll(db *sql.DB) []User {
	var oneUser User
	var users []User
	rows, err := db.Query("select * from user;")
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		err = rows.Scan(&oneUser.id, &oneUser.Name, &oneUser.password, &oneUser.room)
		if err != nil {
			log.Println(err)
		}
		users = append(users, oneUser)
	}
	rows.Close()
	return users
}

//特定のユーザデータを取得
func (user *User) readOne(db *sql.DB) {
	prep, err := db.Prepare("SELECT * FROM user WHERE name = ? LIMIT 1")
	defer prep.Close()
	if err != nil {
		log.Println(err)
	}

	if err = prep.QueryRow(user.Name).Scan(&user.id, &user.Name, &user.password, &user.room); err != nil {
		log.Println(err)
	}
}

func (user *User) idCheck(db *sql.DB) {
	prep, err := db.Prepare("SELECT * FROM user WHERE id = ? LIMIT 1")
	defer prep.Close()
	if err != nil {
		log.Println(err)
	}

	if err = prep.QueryRow(user.id).Scan(&user.id, &user.Name, &user.password, &user.room); err != nil {
		log.Println(err)
	}
}

//新しいデータの追加
func (user *User) insert(db *sql.DB) {
	ins, err := db.Prepare("INSERT INTO user(name,password,room) VALUES(?,?,?)")
	defer ins.Close()
	if err != nil {
		log.Println(err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.password), 12)
	if err != nil {
		log.Println(err)
	}

	log.Println(string(hash))

	_, err = ins.Exec(user.Name, string(hash), user.room)
	if err != nil {
		log.Println(err)
	}
}

//パスワードの変更
func (user *User) update(userOne string, name string, db *sql.DB) {
	var room Room
	room.Name = name
	room.nameCheck(db)
	user.Name = userOne
	user.readOne(db)
	upd, err := db.Prepare("UPDATE user SET room = ? WHERE ( name = ? ) LIMIT 1")
	defer upd.Close()
	if err != nil {
		log.Println(err)
	}
	_, err = upd.Exec(room.Id, user.Name)
	if err != nil {
		log.Println(err)
	}
}

//ユーザデータが存在するかの確認
func (user *User) roomCheck(db *sql.DB) int {
	var oneUser User
	check, err := db.Prepare("SELECT * FROM user WHERE name = ? LIMIT 1")
	defer check.Close()
	if err != nil {
		log.Println(err)
	}
	err = check.QueryRow(user.Name).Scan(&oneUser.id, &oneUser.Name, &oneUser.password, &oneUser.room)
	return oneUser.room
}

func (user *User) userCheck(db *sql.DB) error {
	var oneUser User
	exist, err := db.Prepare("SELECT * FROM user WHERE name = ? LIMIT 1")
	defer exist.Close()
	if err != nil {
		panic(err)
	}
	err = exist.QueryRow(user.Name).Scan(&oneUser.id, &oneUser.Name, &oneUser.password, &oneUser.room)
	log.Println("UserCheck", err)
	return err
}

func (user *User) passCheck(db *sql.DB) error {
	var oneUser User
	exist, err := db.Prepare("SELECT * FROM user WHERE name = ? LIMIT 1")
	defer exist.Close()
	if err != nil {
		log.Println("dbOpen", err)
	}
	err = exist.QueryRow(user.Name).Scan(&oneUser.id, &oneUser.Name, &oneUser.password, &oneUser.room)
	if err == nil {
		err := bcrypt.CompareHashAndPassword([]byte(oneUser.password), []byte(user.password))
		log.Println("PassCheck", err)
		return err
	}
	return err
}

//データの消去z
func (user *User) delete(name string, db *sql.DB) {
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
