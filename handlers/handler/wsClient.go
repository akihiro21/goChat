package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

/*
引用
Anup Kumar Panwar
https://github.com/AnupKumarPanwar/Golang-realtime-chat-rooms
*/

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type sendMes struct {
	Message string `json:"message"`
	Name    string `json:"username"`
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan sendMes
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump() {
	c := s.conn
	defer func() {
		H.unregister <- s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	if err := c.ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.SetFlags(log.Lshortfile)
		log.Println(err)
	}
	err := func(string) error {
		time := time.Now().Add(pongWait)
		if err := c.ws.SetReadDeadline(time); err != nil {
			log.SetFlags(log.Lshortfile)
			log.Println(err)
		}
		return nil
	}
	c.ws.SetPongHandler(err)
	for {
		var msg sendMes
		err := c.ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.SetFlags(log.Lshortfile)
				log.Printf("error: %v", err)
			}
			break
		}

		m := message{msg, s.room}

		H.broadcast <- m
	}
}

// write writes a message with the given message type and payload.
func (c *connection) jsonWrite(payload *sendMes) error {
	if err := c.ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		log.SetFlags(log.Lshortfile)
		log.Println(err)
	}
	return c.ws.WriteJSON(payload)
}

func (c *connection) write(mt int, payload []byte) error {
	if err := c.ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		log.SetFlags(log.Lshortfile)
		log.Println(err)
	}
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, _ := <-c.send:
			if err := c.jsonWrite(&message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request, room string, name string) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.SetFlags(log.Lshortfile)
		log.Println(err)
		return
	}
	c := &connection{send: make(chan sendMes), ws: ws}
	s := subscription{c, room}
	H.register <- s
	go s.writePump()
	go s.readPump()
}
