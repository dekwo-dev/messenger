package main

import (
	"database/sql"
	"log"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := sql.Open("sqlite3", "guestbook.db")
	checkError(err)
	defer db.Close()

    d := newDispatcher()

	ip := getSelfPublicIP()
    url := ip + ":" + port 

    go d.run()

	http.HandleFunc("/comments/post", func (w http.ResponseWriter, 
        r *http.Request) { postComment(db, d, w, r) })
	http.HandleFunc("/comments/get-all", func (w http.ResponseWriter, 
        r *http.Request) { getAllComments(db, w, r) })
    http.HandleFunc("/ws", func (w http.ResponseWriter, 
        r *http.Request) { onNewSubscriber(d, w, r) })

	log.Println("HTTP Server is servering at " + url)

	log.Fatal(http.ListenAndServe(url, nil))
}
