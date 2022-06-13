package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func Open(path string, count uint) *sql.DB {
	db, err := sql.Open("mysql", path)
	if err != nil {
		log.Fatal("open error:", err)
	}

	if err = db.Ping(); err != nil {
		time.Sleep(time.Second * 2)
		count--
		fmt.Printf("retry... count:%v\n", count)
		return Open(path, count)
	}

	fmt.Println("db connected!!")
	return db
}

func ConnectDB() *sql.DB {
	var path string = fmt.Sprintf("%s:%s@tcp(mysql:3306)/%s?charset=utf8&parseTime=true",
		os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_DATABASE"))

	return Open(path, 100)
}
