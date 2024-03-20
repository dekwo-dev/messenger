package test

import (
	"context"
	"log"

	. "dekwo.dev/messager/database"
)

func single_insertion(load uint) {
    conn, err := Pool().Acquire(context.Background())
    if err != nil {
        log.Fatal("Failed to acquire database connection from connection pool")
    }

    tx, err := conn.Begin(context.Background())
    if err != nil {
        log.Println("Failed to start a database transaction. Starting a rollback")
        err := tx.Rollback(context.Background())
        if err != nil {
            log.Fatal("Failed to rollback a database transaction")
        }
    }
    for i := 0; i < int(load); i++ {
        _, err := tx.Exec(context.Background(),
            `INSERT INTO comments (
                commentsid,
                commentsauthor,
                commentscontent,
                commentstime
            ) VALUES (DEFAULT, 'Author $1', 'Comment $2', DEFAULT)`,
            i, i,
        )
        if err != nil {
            log.Println("Failed to start a database transaction. Starting a rollback")
            err := tx.Rollback(context.Background())
            if err != nil {
                log.Fatal("Failed to rollback a database transaction")
            }
        }
    }
    err = tx.Commit(context.Background())
    if err != nil {
        log.Println("Failed to start a database transaction. Starting a rollback")
        err := tx.Rollback(context.Background())
        if err != nil {
            log.Fatal("Failed to rollback a database transaction")
        }
    }
}
