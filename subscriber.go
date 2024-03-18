package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"dekwo.dev/messager/pb"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  8,
	WriteBufferSize: 32,
	CheckOrigin: func(r *http.Request) bool {
		o := r.Header.Get("Origin")
		local := strings.HasPrefix(o, "http://localhost") ||
			strings.HasPrefix(o, "http://127.0.0.1")
        remote := strings.HasPrefix(o, "https://dekr0.dev")
        if prod() {
            return remote
        } else {
            return local
        }
	},
}

type Subscriber struct {
	conn  *websocket.Conn
	addr  string
	event chan *pb.Event      // a subscriber <-> dispatcher
	unsub chan *Subscriber // a subscriber <-> dispatcher
}

func (sub *Subscriber) close() {
    close(sub.event)
    sub.unsub <- sub
    sub.conn.Close()
}

func (sub *Subscriber) read() {
    const f = "Subscriber.read"
    const file = "subscriber.go"

    defer sub.close()

	sub.conn.SetReadLimit(8);

	for {
		_, _, err := sub.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
                info(50, file, f, fmt.Sprintf("Failed to close connection from %s",
                    sub.addr), err)
			}
			break
		}
	}

	info(50, file, f, fmt.Sprintf("Connection from %s closed", sub.addr), nil)
}

func (sub *Subscriber) notify() {
    const f = "Subscriber.notify"
    const file = "subscriber.go"

    defer sub.close()

	info(50, file, f, fmt.Sprintf("Subscriber from %s is listening for next event",
        sub.addr), nil)

    attempt := 5

    for {
        select {
		case event := <-sub.event:
			sub.conn.SetWriteDeadline(time.Now().Add(writeWait))

            info(30, file, f, fmt.Sprintf("Subscriber from %s receive event: %s",
                    sub.addr, event.String()), nil)

            if err := sub.conn.WriteMessage(
                websocket.TextMessage, 
                []byte(event.String()),
            ); err != nil {
                info(50, file, f, fmt.Sprintf("Failed to write payload to Subscriber from %s. Remaining failed write attempt: %d",
                        sub.addr, attempt - 1), err)
                attempt--
                if attempt <= 0 {
                    info(50, file, f, fmt.Sprintf("Failed write attempt limit reaches for Subscriber from %s. Forcefully close connection", 
                            sub.addr), nil)
                    return
                }
			}
		}
	}
}

func ws(d *Dispatcher, w http.ResponseWriter, r *http.Request) {
    const f = "ws"
    const file = "subscriber.go"

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
        info(50, file, f, "Failed to accept a WS connection", err)
		return
	}

	sub := &Subscriber{
		conn:  conn,
		addr:  conn.RemoteAddr().String(),
		event: make(chan *pb.Event),
		unsub: d.unsub,
	}

	d.sub <- sub

	go sub.notify()
	go sub.read()
}
