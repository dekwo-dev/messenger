package worker

import (
    "context"
    "strconv"
    "strings"

    . "dekwo.dev/messager/database"
    . "dekwo.dev/messager/logger"
)

type Notifier struct {
    enqueue chan int
}

func NewNotifier() *Notifier {
    return &Notifier{
        make(chan int),
    } 
}

func (n *Notifier) Run() {
    const f = "Notifier.Run"
    const file = "worker/notifier.go"

    defer func() {
        close(n.enqueue)
    }()

    conn, err := Pool().Acquire(context.Background())
    if err != nil {
        Info(50, file, f, "Failed to acquire a pgx Connection from a pool", err)
        return
    }
    defer conn.Release()

    _, err = conn.Exec(context.Background(), "LISTEN enqueue_comments_changes")
    if err != nil {
        Info(50, file, f, "Failed to execute the LISTENING statement", err)
        return
    }

    for {
        Info(20, file, f, "Notifier is waiting from the next notification", nil)
        in, err := conn.Conn().WaitForNotification(context.Background())
        if err != nil {
            Info(50, file, f, "Error when receiving the incoming notification", err)
        }
        load, err := strconv.ParseInt(strings.Split(in.Payload, ":")[1], 10, 8)
        if err != nil {
            Info(50, file, f, "Error when parsing the notification payload", err)
        }

        n.enqueue <- int(load)
    }
}
