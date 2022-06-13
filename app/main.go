package main

import (
	"log"
	"net/http"
	"os"

	"github.com/akihiro21/goChat/handlers/database"
	"github.com/akihiro21/goChat/handlers/handler"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	myDb := database.ConnectDB()
	defer myDb.Close()
	mux := http.NewServeMux()
	dir, _ := os.Getwd()
	port := "8080"
	log.Printf("Server listening on http://localhost:%s/", port)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(dir+"/web/static/"))))

	handler.Init(mux, myDb)
	log.Print(http.ListenAndServe(":"+port, mux))
}
