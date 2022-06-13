module github.com/akihiro21/goChat

go 1.18

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gorilla/sessions v1.2.1
	github.com/gorilla/websocket v1.5.0
	github.com/pkg/errors v0.9.1
	golang.org/x/crypto v0.0.0-20220516162934-403b01795ae8
)

replace github.com/akihiro21/goChat/handlers => ./handlers

require github.com/gorilla/securecookie v1.1.1 // indirect
