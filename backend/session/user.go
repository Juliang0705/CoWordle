package session

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var upgrader = websocket.Upgrader{}

type User struct {
	controller       *Controller
	session          *Session
	conn             *websocket.Conn
	outboundResponse chan ServerResponse
}

func (u *User) DoReadLoop() {
	defer func() {
		if u.session != nil {
			u.session.unregisterUser <- u
		}
		u.conn.Close()
	}()
	u.conn.SetReadLimit(maxMessageSize)
	u.conn.SetReadDeadline(time.Now().Add(pongWait))
	u.conn.SetPongHandler(func(string) error { u.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var command UserCommand = UserCommand{}
		err := u.conn.ReadJSON(&command)
		log.Printf("----------------------------------------\n")
		log.Printf("reading: %+v", command)
		log.Printf("----------------------------------------\n")
		if err != nil {
			log.Printf("error: %v", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return
		}
		err = u.controller.ProcessCommand(command, u)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
	}
}

func (u *User) DoWriteLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		u.conn.Close()
		if u.session != nil {
			u.session.unregisterUser <- u
		}
	}()
	for {
		select {
		case response, ok := <-u.outboundResponse:
			u.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The session closed the channel.
				u.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := u.conn.WriteJSON(response)
			log.Printf("writing: %+v", response)
			if err != nil {
				log.Printf("error: %v", err)
				return
			}
		case <-ticker.C:
			u.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := u.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeUserWebsocket(controller *Controller, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	user := &User{controller: controller, conn: conn, outboundResponse: make(chan ServerResponse, maxMessageSize)}
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go user.DoReadLoop()
	go user.DoWriteLoop()
}
