package worker

import (
    "context"
    "fmt"
    "strings"
    "time"

    . "dekwo.dev/messager/database"
    . "dekwo.dev/messager/logger"
    "dekwo.dev/messager/pb"
)

type Worker struct {
    done    chan *Worker
    payload []*pb.Event
}

func newWorker(d *Dispatcher) *Worker {
    return &Worker{ 
        d.done,
        []*pb.Event{},
    }
}

func (w *Worker) run(load int) {
    const f = "Worker.run"
    const file = "worker/worker.go"

    conn, err := Pool().Acquire(context.Background())
    if err != nil {
        Info(50, file, f, "Failed to acquire a pgx Connection from a pool", err)
        return
    }
    defer conn.Release()

    tx, err := conn.Begin(context.Background())
    if err != nil {
        Info(50, file, f, "Failed to start a transaction", err)
        return
    }

    rows, err := tx.Query(context.Background(),
        `DELETE FROM comments_changes_queue
        WHERE queueid in (
            SELECT queueid FROM comments_changes_queue
            FOR UPDATE SKIP LOCKED
            LIMIT $1
        )
        RETURNING operation, commentsid, commentsauthor, commentscontent, commentstime`,
        load,
    )
    if err != nil {
        Info(50, file, f, "Failed to execute query", err)
        if err := tx.Rollback(context.Background()); err != nil {
            Info(50, file, f, "Failed to perform transaction rollback after query failed",
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

        if err = rows.Scan(
            &operation,
            &commentsid,
            &commentsauthor,
            &commentscontent,
            &commentstime,
        ); err != nil {
            Info(50, file, f, "Failed to scan the next queue row", err)
            if err := tx.Rollback(context.Background()); err != nil {
                Info(50, file, f, "Failed to perform transaction rollback after query failed",
                    nil)
            }
            return
        }

        event := &pb.DBChangeEvent {
            DbEventType: pb.DBEventType(pb.DBEventType_value[fmt.Sprintf("DBEVENTTYPE_%s", operation)]),
            Comment: &pb.Comment { Id: uint32(commentsid) },
        }

        wrapper := &pb.Event { 
            EventOneof: &pb.Event_DbChangeEvent { DbChangeEvent: event, },
        }

        insert := strings.Compare(operation, "INSERT") == 0
        update := strings.Compare(operation, "UPDATE") == 0        
        if insert || update {
            event.Comment.Author  = commentsauthor
            event.Comment.Content = commentscontent 
            event.Comment.Time    = commentstime.UnixMilli()  
        }

        w.payload = append(w.payload, wrapper)
    }

    if err = tx.Commit(context.Background()); err != nil {
        Info(50, f, file, "Failed to commit the work in the work queue", nil)
            if err := tx.Rollback(context.Background()); err != nil {
            Info(50, file, f, "Failed to perform transaction rollback after work commit failed",
                nil)
        }
        return
    }

    w.done <- w
}
