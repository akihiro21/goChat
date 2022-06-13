module github.com/akihiro21/goChat

go 1.18

require (
	github.com/akihiro21/goChat/handlers v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.6.0
)

require (
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/gorilla/sessions v1.2.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
)

replace github.com/akihiro21/goChat/handlers => ./handlers
