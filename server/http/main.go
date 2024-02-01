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

	_ "github.com/mattn/go-sqlite3"
)

const debug = true
var db *sql.DB
var connection net.Conn

func checkError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func postCommentHandler(w http.ResponseWriter, r *http.Request) {
    // dummy, comments
    nextCommentsId:= getNextCommentsId()

    transaction, err := db.Begin()
    checkError(err)

    statement, err := transaction.Prepare(`INSERT INTO comments(commentsId, 
    commentsAuthor, commentsContent, commentsTime) values(?, ?, ?, ?)`)
    checkError(err)
    defer statement.Close()

    _, err = statement.Exec(nextCommentsId + 1, "user_1", "This is a comment", 
    time.Now())
    checkError(err) 

    err = transaction.Commit()
    checkError(err)
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
        checkError(err)
        defer rows.Close()

        for rows.Next() {
            var (
                id int
                author string
                content string
                timestamp string
            )

            err := rows.Scan(&id, &author, &content, &timestamp)
            checkError(err)
            fmt.Println(id, author, content, timestamp)
        }

    } 
}

func ping(w http.ResponseWriter, r *http.Request) {
    onDBChange()
}

// testing purpose
func onDBChange() {
}

func getNextCommentsId() int {
    statement := `SELECT MAX(commentsId) as nextCommentsId FROM comments`
    
    rows, err := db.Query(statement)
    checkError(err)
    defer rows.Close()

    if rows.Next() {
        var id int
        err = rows.Scan(&id)
        checkError(err)

        return id
    }

    return -1
}

func main() {
    var ipBuilder strings.Builder
    var buffer []byte
    ipBuilder.WriteByte(':')

    var err error
    port := "8080"

    db, err = sql.Open("sqlite3", "guestbook.db")
    checkError(err)
    defer db.Close()
    
    if !debug{
        ipBuilder.Reset()
        var response *http.Response
        response, err = http.Get("https://api.ipify.org?format=text") 
        checkError(err)

        defer response.Body.Close()
        buffer, err = io.ReadAll(response.Body)
        checkError(err)
        
        ipBuilder.Write(buffer)
        ipBuilder.WriteByte(':')
    }

    ip := ipBuilder.String()
    ipBuilder.Reset()

    http.HandleFunc("/comments/post", postCommentHandler)
    http.HandleFunc("/comments/get-all", getAllCommentHandler)
    http.HandleFunc("/comments/get-comment", getCommentById)
    http.HandleFunc("/comments/ping", ping) 
    log.Println("HTTP Server is servering at " + ip + port)
    log.Fatal(http.ListenAndServe(ip + port, nil))
}
