package main

import (
	"log"
	"net/http"
)

func main() {
    n := notifier()
	d := dispatcher(n)

	ip := getSelfPublicIP()
	url := ip + ":" + port

	go d.run()
    go n.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter,
		r *http.Request) {
		ws(d, w, r)
	})

	log.Fatal(http.ListenAndServe(url, nil))
}
