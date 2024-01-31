package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB


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
    var err error
    db, err = sql.Open("sqlite3", "guestbook.db")
    checkError(err)

    defer db.Close()

    http.HandleFunc("/comments/post", postCommentHandler)
    http.HandleFunc("/comments/get-all", getAllCommentHandler)
    http.HandleFunc("/comments/get-comment", getCommentById)
    fmt.Println("Launching http server")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
