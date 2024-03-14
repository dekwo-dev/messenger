package main

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool = pool()

type Notifier struct {
    enqueue chan int
}

type Worker struct {
    done    chan *Worker
    payload []Event 
}

func pool() func() *pgxpool.Pool {
    const f = "pool"
    const file = "db.go"

    var pool *pgxpool.Pool

    return func() *pgxpool.Pool {
        if pool != nil {
            return pool
        }

        url := GetEnv("POSTGRES_URL")
        if strings.Compare(url, "") == 0 {
            fatal(50, file, f, "Postgres connection url is required in .env", nil) 
        }

        cfg, e := pgxpool.ParseConfig(url) 
        if e != nil {
            fatal(50, file, f, "Failed to parse Postgres connection url", e)
        }

        if pool, e = pgxpool.NewWithConfig(context.Background(), cfg); e != nil {
            fatal(50, file, f, "Failed to create a new pgxpool.Pool", e)
        }

        return pool
    }
}

func notifier() *Notifier {
    return &Notifier{
        make(chan int),
    } 
}

func worker() *Worker {
    return &Worker{
        make(chan *Worker),
        make([]Event, 1),
    }
}

func (n *Notifier) close() {
    close(n.enqueue)
}

func (n *Notifier) run() {
    const f = "Notifier.run"
    const file ="db.go"

    defer n.close()

    conn, e := Pool().Acquire(context.Background())
    if e != nil {
        info(50, file, f, "Failed to acquire a pgx Connection from a pool", e)
        return
    }
    defer conn.Release()

    _, e = conn.Exec(context.Background(), "LISTEN enqueue_comments_changes")
    if e != nil {
        info(50, file, f, "Failed to execute the LISTENING statement", e)
        return
    }

    for {
        info(30, file, f, "Notifier is waiting from the next notification", nil)
        in, e := conn.Conn().WaitForNotification(context.Background())
        if e != nil {
            info(50, file, f, "Error when receiving the incoming notification", e)
        }

        load, e := strconv.ParseInt(strings.Split(in.Payload, ":")[1], 10, 1)
        if e != nil {
            info(50, file, f, "Error when parsing the notification payload", e)
        }

        n.enqueue <- int(load)
    }
}

func (w *Worker) close() {
    close(w.done)
}

func (w *Worker) run(load int) {
    const f = "Worker.run"
    const file = "db.go"

    defer w.close()

    conn, e := Pool().Acquire(context.Background())
    if e != nil {
        info(50, file, f, "Failed to acquire a pgx Connection from a pool", e)
        return
    }
    defer conn.Release()

    tx, e := conn.Begin(context.Background())
    if e != nil {
        info(50, file, f, "Failed to start a transaction", e)
        return
    }

    rows, e := tx.Query(context.Background(),
        `DELETE FROM comments_changes_queue
        WHERE queueid in (
            SELECT queueid FROM comments_changes_queue
            FOR UPDATE SKIP LOCKED
            LIMIT $1
        )
        RETURNING operation, commentsid, commentsauthor, commentscontent, commentsauthor`,
    )
    if e != nil {
        info(50, file, f, "Failed to execute query", e)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var (
            operation string
            commentsid int
            commentsauthor string
            commentscontent string
            commentstime time.Time
        )

        if e = rows.Scan(
            &operation,
            &commentsid,
            &commentsauthor,
            &commentscontent,
            &commentstime,
        ); e != nil {
            info(50, file, f, "Failed to scan the next queue row", e)
            return
        }
    }

    w.done <- w
}
