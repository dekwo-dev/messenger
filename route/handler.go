package route

import (
    "net/http"

    . "dekwo.dev/messager/worker"
    . "dekwo.dev/messager/logger"
)

func ws(d *Dispatcher, w http.ResponseWriter, r *http.Request) {
    const f = "ws"
    const file = "route/handler.go"

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
        Info(50, file, f, "Failed to accept a WS connection", err)
		return
	}

	sub := NewSubscriber(conn, d)
	d.Sub <- sub

	go sub.Notify()
	go sub.Read()
}

