package main

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"dekwo.dev/messager/pb"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool = pool()

type Notifier struct {
    enqueue chan int
}

type Worker struct {
    done    chan *Worker
    payload []*pb.Event
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

func worker(d *Dispatcher) *Worker {
    return &Worker{ 
        d.done,
        []*pb.Event{},
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

        log.Println(in.Payload)
        load, e := strconv.ParseInt(strings.Split(in.Payload, ":")[1], 10, 8)
        if e != nil {
            info(50, file, f, "Error when parsing the notification payload", e)
        }

        n.enqueue <- int(load)
    }
}

func (w *Worker) run(load int) {
    const f = "Worker.run"
    const file = "db.go"

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
        RETURNING operation, commentsid, commentsauthor, commentscontent, commentstime`,
        load,
    )
    if e != nil {
        info(50, file, f, "Failed to execute query", e)
        if e = tx.Rollback(context.Background()); e != nil {
            info(50, file, f,
                "Failed to perform transaction rollback after query failed",
                nil)
        }
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
            if e = tx.Rollback(context.Background()); e != nil {
                info(50, file, f,
                    "Failed to perform transaction rollback after query failed",
                    nil)
            }
            return
        }

        event := &pb.DBChangeEvent {
            DbEventType: pb.DBEventType(pb.DBEventType_value[operation]),
            Comment: &pb.Comment { Id: uint32(commentsid) },
        }

        wrapper := &pb.Event { 
            EventOneof: &pb.Event_DbChangeEvent { DbChangeEvent: event, },
        }

        insert := strings.Compare(operation, "INSERT") == 0
        update := strings.Compare(operation, "INSERT") == 0        
        if insert || update {
            event.Comment.Author  = commentsauthor
            event.Comment.Content = commentscontent 
            event.Comment.Time    = commentstime.UnixMilli()  
        }

        w.payload = append(w.payload, wrapper)
    }

    if e = tx.Commit(context.Background()); e != nil {
        info(50, f, file,
            "Failed to commit the work in the work queue",
            nil)
        if e = tx.Rollback(context.Background()); e != nil {
            info(50, file, f,
                "Failed to perform transaction rollback after work commit failed",
                nil)
        }
        return
    }

    w.done <- w
}

