package main

import (
	"log"
	"net/http"
	"os"

	"github.com/akihiro21/goChat/packages"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	packages.SessionInit()

	myDb := packages.ConnectDB()
	defer myDb.Close()

	packages.Init(myDb)

	mux := http.NewServeMux()
	dir, _ := os.Getwd()
	port := "8080"
	log.Printf("Server listening on http://localhost:%s/", port)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(dir+"/web/static/"))))

	packages.HandleInit(mux)
	log.Print(http.ListenAndServe(":"+port, mux))
}
