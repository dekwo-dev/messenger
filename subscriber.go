package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
    ReadBufferSize: 8,
    WriteBufferSize: 32,
    CheckOrigin: func (r *http.Request) bool {
        o := r.Header.Get("Origin")
        log.Printf("CheckOrigin (INFO): %v", o)
        isLocalHost := strings.HasPrefix(o, "http://localhost") || 
            strings.HasPrefix(o, "http://127.0.0.1")
        if isLocalHost {
            return true
        } else {
            return false
        }
    },
}

type Subscriber struct {
	conn     *websocket.Conn
    addr string
	event chan []byte // a subscriber <-> dispatcher
    unsub chan *Subscriber // a subscriber <-> dispatcher
}

func (sub *Subscriber) readSub() {
	// This function will run after blocking read loop is broke, which indicate
	// a subscriber is unsubscribe
    const f = "readSub"

    defer func() {
        sub.unsub <- sub
        sub.conn.Close()
    }()

	sub.conn.SetReadLimit(8)
    log.Printf(
        "%v at %v (INFO): Waiting for subscriber to close connection\n", 
        f, sub.addr,
    )
	for {
		// ignore subscriber msg
		_, _, err := sub.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("%v at %v (ERROR): %v\n", f, sub.addr, err)
			}
			break
		}
	}
    log.Printf("%v at %v (INFO): Connection is closed\n", f, sub.addr)
}

func (sub *Subscriber) notifySub() {
    defer func() {
        sub.unsub <- sub
	    sub.conn.Close()
    }()

    const f = "writeSub"

    log.Printf("%v at %v (INFO): Waiting for new DB change\n", f, sub.addr)
	for {
		select {
		case event := <-sub.event:
            sub.conn.SetWriteDeadline(time.Now().Add(writeWait))
            log.Printf("%v at %v (INFO): event - %v\n", 
                f, sub.addr, string(event))
            err := sub.conn.WriteMessage(websocket.TextMessage, event)
			if err != nil {
                log.Printf("%v at %v (ERROR): Failed to notify\n", f, sub.addr)
                log.Println(err)
                log.Printf("%v at %v (ERROR): Closing connection\n", f, sub.addr)
                return
			}
		}
	}
}

func onNewSubscriber(d *Dispatcher, w http.ResponseWriter, r *http.Request) {
    const f = "onNewSubscriber"

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("%v (ERROR): Failed to accept a subscriber", f)
        log.Println(err)
        return
    }
    
    sub := &Subscriber{
        conn: conn,
        addr: conn.RemoteAddr().String(),
        event: make(chan []byte),
        unsub: d.unsub,
    }
    
    d.sub <- sub

    log.Printf("%v (INFO): New subscriber at %v\n", f, sub.addr)

    go sub.notifySub()
    go sub.readSub()
}
