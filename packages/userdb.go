package packages

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id       int
	name     string
	password string
	room     int
}

//全てのユーザデータを取得
func (user *User) ReadAll(db *sql.DB) []User {
	var oneUser User
	var users []User
	rows, err := db.Query("select * from user;")
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		err = rows.Scan(&oneUser.id, &oneUser.name, &oneUser.password, &oneUser.room)
		if err != nil {
			log.Println(err)
		}
		users = append(users, oneUser)
	}
	rows.Close()
	return users
}

//特定のユーザデータを取得
func (user *User) ReadOne(db *sql.DB) {
	prep, err := db.Prepare("SELECT * FROM user WHERE name = ? LIMIT 1")
	defer prep.Close()
	if err != nil {
		log.Println(err)
	}

	if err = prep.QueryRow(user.name).Scan(&user.id, &user.name, &user.password, &user.room); err != nil {
		log.Println(err)
	}
}

//新しいデータの追加
func (user *User) Insert(db *sql.DB) {
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

	_, err = ins.Exec(user.name, string(hash), user.room)
	if err != nil {
		log.Println(err)
	}
}

//パスワードの変更
func (user *User) Update(id int, db *sql.DB) {
	upd, err := db.Prepare("UPDATE user SET room = ? WHERE ( name = ? ) LIMIT 1")
	defer upd.Close()
	if err != nil {
		log.Println(err)
	}
	_, err = upd.Exec(id, user.name)
	if err != nil {
		log.Println(err)
	}
}

//ユーザデータが存在するかの確認
func (user *User) Exist(db *sql.DB) error {
	var oneUser User
	exist, err := db.Prepare("SELECT * FROM user WHERE name = ? and password = ? LIMIT 1")
	defer exist.Close()
	if err != nil {
		log.Println(err)
	}
	err = exist.QueryRow(user.name, user.password).Scan(&oneUser.id, &oneUser.name, &oneUser.password, &oneUser.room)
	return err
}

func (user *User) UserCheck(db *sql.DB) error {
	var oneUser User
	exist, err := db.Prepare("SELECT * FROM user WHERE name = ? LIMIT 1")
	defer exist.Close()
	if err != nil {
		panic(err)
	}
	err = exist.QueryRow(user.name).Scan(&oneUser.id, &oneUser.name, &oneUser.password, &oneUser.room)
	log.Println("UserCheck", err)
	return err
}

func (user *User) PassCheck(db *sql.DB) error {
	var oneUser User
	exist, err := db.Prepare("SELECT * FROM user WHERE name = ? LIMIT 1")
	defer exist.Close()
	if err != nil {
		log.Println("dbOpen", err)
	}
	err = exist.QueryRow(user.name).Scan(&oneUser.id, &oneUser.name, &oneUser.password, &oneUser.room)
	if err == nil {
		err := bcrypt.CompareHashAndPassword([]byte(oneUser.password), []byte(user.password))
		log.Println("PassCheck", err)
		return err
	}
	return err
}

//データの消去z
func (user *User) Delete(db *sql.DB) {
	delete, err := db.Prepare("DELETE FROM user WHERE name = ? ")
	defer delete.Close()
	if err != nil {
		log.Println(err)
	}
	_, err = delete.Exec(user.name, user.password)
	defer delete.Close()
	if err != nil {
		log.Println(err)
	}
}
