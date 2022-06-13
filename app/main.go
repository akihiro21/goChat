package main

import (
	"log"
	"net/http"
	"os"

	"github.com/akihiro21/goChat/handlers"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	handlers.SessionInit()

	myDb := handlers.ConnectDB()
	defer myDb.Close()

	handlers.Init(myDb)

	mux := http.NewServeMux()
	dir, _ := os.Getwd()
	port := "8080"
	log.Printf("Server listening on http://localhost:%s/", port)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(dir+"/web/static/"))))

	handlers.HandleInit(mux)
	log.Print(http.ListenAndServe(":"+port, mux))
}
