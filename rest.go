package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func postComment(db *sql.DB, d *Dispatcher, w http.ResponseWriter, 
    r *http.Request) {

    const f = "postCommnet"

	id, err := latestCommentsId(db)
    if err != nil {
        log.Printf("%v (ERROR): Failed to fetch next commentsId\n", f)
        log.Println(err)
        return
    }

	t, err := db.Begin()
	if err != nil {
		log.Printf("%v (ERROR): Failed to start a DB transaction\n", f)
		log.Println(err)
		return
	}

	stmt, err := t.Prepare(`INSERT INTO comments(commentsId, 
    commentsAuthor, commentsContent, commentsTime) values(?, ?, ?, ?)`)
	if err != nil {
		log.Printf("%v (ERROR): Failed to prepare insert stmt\n", f)
		log.Println(err)
		return
	}
	defer stmt.Close()

	// dummy comments
	_, err = stmt.Exec(id + 1, "user_1", "This is a comment",
		time.Now())
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
    

    // http server -> Dispather -> every subscriber
    d.event <- []byte(NewComment)
}

func getAllComments(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var builder strings.Builder

    const f = "getAllComments"

	buf := make([]byte, 1024, 1024)
	for {
		_, err := r.Body.Read(buf)
		if err == nil {
			builder.Write(buf)
		} else if err == io.EOF {
			builder.Write(buf)
			break
		} else {
			log.Printf(`%v (ERROR): Error when reading request body %v\n`, f, err)
			return
		}
	}
	w.Write([]byte(builder.String()))
}

// for testing purpose
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
				timestamp string
			)

			err := rows.Scan(&id, &author, &content, &timestamp)
			checkError(err)
			fmt.Println(id, author, content, timestamp)
		}

	}
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
