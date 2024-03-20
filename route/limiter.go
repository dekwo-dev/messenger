package route

import (
	"net/http"
	"strings"

    . "dekwo.dev/messager/env"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
	ReadBufferSize:  8,
	WriteBufferSize: 32,
	CheckOrigin: func(r *http.Request) bool {
		o := r.Header.Get("Origin")
		local := strings.HasPrefix(o, "http://localhost") ||
			strings.HasPrefix(o, "http://127.0.0.1")
        remote := strings.HasPrefix(o, "https://dekr0.dev")
        if Prod() {
            return remote
        } else {
            return local
        }
	},
}

func limiter() {

}
