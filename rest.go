package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Comment struct {
	Id        int
	Author    string
	Content   string
	Timestamp int64 
}

func p(reader io.ReadCloser) {
    var builder strings.Builder
    buf := make([]byte, 128)
    for {
        n, err := reader.Read(buf)
        if err != nil {
            if err == io.EOF {
                builder.Write(buf[:n])
                log.Printf("p (INFO): %v", builder.String())
                builder.Reset()
                return
            } else {
                log.Printf("p (ERROR): %v", err)
                return
            }
        }
        builder.Write(buf[:n])
    }
}

func postComment(db *sql.DB, d *Dispatcher, w http.ResponseWriter,
	r *http.Request) {
	const f = "postCommnet"

    if err := r.ParseForm(); err != nil {
        log.Printf("%v (ERROR): Failed to parse form. Reason: %v\n", f, err)
    }
    author := r.FormValue("author")
    comment := r.FormValue("comment")
    log.Printf("%v (INFO): author = %v; comment = %v", f, author, comment)

	id, err := latestCommentsId(db)
	if err != nil {
		log.Printf("%v (ERROR): Failed to fetch next commentsId. Reason: %v\n",
			f, err)
		return
	}

	t, err := db.Begin()
	if err != nil {
		log.Printf("%v (ERROR): Failed to start a DB transaction. Reason: %v\n",
			f, err)
		return
	}

	stmt, err := t.Prepare(`INSERT INTO comments(commentsId, 
    commentsAuthor, commentsContent, commentsTime) values(?, ?, ?, ?)`)
	if err != nil {
		log.Printf("%v (ERROR): Failed to prepare insert stmt. Reason: %v\n",
			f, err)
		return
	}
	defer stmt.Close()

    now := time.Now().UnixMilli()

	_, err = stmt.Exec(id+1, author, comment, now)
	if err != nil {
		log.Printf("%v (ERROR): Failed to execute insert stmt\n", f)
		log.Println(err)
		return
	}

	err = t.Commit()
	if err != nil {
		log.Printf("%v (ERROR): Failed to commit DB transaction\n", f)
		log.Println(err)
		return
	}

    e := &DBChangeEvent{
        "DBChangeEvent",
        NewComment,
        Comment{
            id + 1,
            author,
            comment,
            now,
        },
    };

	// http server -> Dispather -> every subscriber
	d.event <- e

    w.WriteHeader(200);
}

func getAllComments(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	const f = "getAllComments"
	rows, err := db.Query("SELECT * FROM comments ORDER BY commentsTime DESC")
	if err != nil {
		log.Printf("%v (ERROR): Failed to execute query. Reason: %v\n", f, err)
		return
	}
	defer rows.Close()
	comments := []Comment{}
	for rows.Next() {
		var (
			id        int
			author    string
			content   string
			timestamp int64 
		)

		if err := rows.Scan(&id, &author, &content, &timestamp); err != nil {
			log.Printf("%v (ERROR): Failed to scan row. Reason: %v\n", f, err)
			continue
		}

		comments = append(comments, Comment{id, author, content, timestamp})
	}
	b, err := json.Marshal(comments)

	w.Write(b)
}

func latestCommentsId(db *sql.DB) (int, error) {
	stmt := `SELECT 1 FROM comments`
	rows, err := db.Query(stmt)
	if err != nil {
		return -2, err
	}

	if !rows.Next() {
		return -1, nil
	}
	rows.Close()

	stmt = `SELECT MAX(commentsId) as latestCommentsId FROM comments`

	rows, err = db.Query(stmt)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err == sql.ErrNoRows {
			return -1, nil
		} else if err != nil {
			return -2, err
		}

		return id, nil
	}

	return -1, nil
}

func getCommentById(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Has("commentsId") {
		commentsId := r.URL.Query().Get("commentsId")

		rows, err := db.Query("SELECT * FROM comments WHERE commentsId = ?",
			commentsId)
		checkError(err)
		defer rows.Close()

		for rows.Next() {
			var (
				id        int
				author    string
				content   string
				timestamp int64 
			)

			err := rows.Scan(&id, &author, &content, &timestamp)
			checkError(err)
			fmt.Println(id, author, content, timestamp)
		}

	}
}
