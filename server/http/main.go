package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"example.com/guestbook/config"
	"example.com/guestbook/utils"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var connection net.Conn

func postCommentHandler(w http.ResponseWriter, r *http.Request) {
    // dummy, comments
    nextCommentsId:= getNextCommentsId()

    transaction, err := db.Begin()
    utils.CheckError(err)

    statement, err := transaction.Prepare(`INSERT INTO comments(commentsId, 
    commentsAuthor, commentsContent, commentsTime) values(?, ?, ?, ?)`)
    utils.CheckError(err)
    defer statement.Close()

    _, err = statement.Exec(nextCommentsId + 1, "user_1", "This is a comment", 
    time.Now())
    utils.CheckError(err)

    err = transaction.Commit()
    utils.CheckError(err)
}

func getAllCommentHandler(w http.ResponseWriter, r *http.Request) {
    var builder strings.Builder
    fmt.Println(r.Header.Get("Content-Type"))
    buffer := make([]byte, 1024, 1024)
    for {
        _, err := r.Body.Read(buffer) 
        if err == nil || err == io.EOF {
            builder.Write(buffer)
            break
        } else {
            log.Fatalf(`GetAllCommentHandler: Error when reading request 
            body %v\n`, err)
        }
        builder.Write(buffer)
    }
    fmt.Println(builder.String())
}

// for testing purpose
func getCommentById(w http.ResponseWriter, r *http.Request) {
    if r.URL.Query().Has("commentsId") {
        commentsId := r.URL.Query().Get("commentsId")

        rows, err := db.Query("SELECT * FROM comments WHERE commentsId = ?", 
        commentsId)
        utils.CheckError(err)
        defer rows.Close()

        for rows.Next() {
            var (
                id int
                author string
                content string
                timestamp string
            )

            err := rows.Scan(&id, &author, &content, &timestamp)
            utils.CheckError(err)
            fmt.Println(id, author, content, timestamp)
        }

    } 
}

func ping(w http.ResponseWriter, r *http.Request) {
}


func getNextCommentsId() int {
    statement := `SELECT MAX(commentsId) as nextCommentsId FROM comments`
    
    rows, err := db.Query(statement)
    utils.CheckError(err)
    defer rows.Close()

    if rows.Next() {
        var id int
        err = rows.Scan(&id)
        utils.CheckError(err)

        return id
    }

    return -1
}

func dialEventEmitter() {
    var err error
    
    url := "localhost" + ":" + config.SocketPort

    if config.EnableRemoteSocket {
        url = config.RemoteIP + ":" + config.SocketPort
    }
    connection, err = net.Dial("tcp", url)
    utils.CheckError(err)
}

func main() {
    var err error

    db, err = sql.Open("sqlite3", "guestbook.db")
    utils.CheckError(err)
    defer db.Close()
    
    ip := utils.GetSelfPublicIP()
    
    url := ip + config.HTTPPort

    dialEventEmitter()

    http.HandleFunc("/comments/post", postCommentHandler)
    http.HandleFunc("/comments/get-all", getAllCommentHandler)
    http.HandleFunc("/comments/get-comment", getCommentById)
    http.HandleFunc("/comments/ping", ping) 
    log.Println("HTTP Server is servering at " + url)
    log.Fatal(http.ListenAndServe(url, nil))
}
