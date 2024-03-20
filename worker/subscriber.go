package worker 

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	"dekwo.dev/messager/pb"
    "dekwo.dev/messager/logger"
)

type Subscriber struct {
	conn  *websocket.Conn
	addr  string
	event chan *pb.Event      // a subscriber <- dispatcher
	unsub chan *Subscriber // a subscriber -> dispatcher
    loop  chan *E
}

func NewSubscriber(c *websocket.Conn, d *Dispatcher) *Subscriber {
    return &Subscriber{
		conn:  c,
		addr:  c.RemoteAddr().String(),
		event: make(chan *pb.Event),
        loop:  make(chan *E),
		unsub: d.unsub,
	}
}

func (sub *Subscriber) Read() {
    const f = "Subscriber.Read"
    const file = "worker/subscriber.go"

    defer func() {
        logger.Info(20, file, f, fmt.Sprintf("Subscriber from %s signals self notify work to stop notifying",
            sub.addr), nil)
        sub.loop <- &E{}
    }()

	sub.conn.SetReadLimit(8);

	for {
		_, _, err := sub.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
                logger.Info(40, file, f, fmt.Sprintf("Connection Close from %s", 
                    sub.addr), err)
			} else {
                logger.Info(50, file, f, fmt.Sprintf("Failed to ignore message from %s", 
                    sub.addr), err)
            }
            break
		}
	}

	logger.Info(20, file, f, fmt.Sprintf("Connection from %s closed", sub.addr), nil)
}

func (sub *Subscriber) Notify() {
    const f = "Subscriber.Notify"
    const file = "worker/subscriber.go"

    const writeWait = 10 * time.Second

    defer func() {
        close(sub.event)
        close(sub.loop)
        sub.unsub <- sub
        sub.conn.Close()
        sub.conn = nil
        logger.Info(20, file, f, fmt.Sprintf("Subscriber from %s: Cleanup finish", 
            sub.addr), nil)
    }()

	logger.Info(20, file, f, fmt.Sprintf("Subscriber from %s is listening for next event",
        sub.addr), nil)

    attempt := 5

    for {
        select {
        case <- sub.loop:
            logger.Info(20, file, f, fmt.Sprintf("Stopped notifying work for %s", 
                sub.addr), nil)
            return
		case event := <-sub.event:
			sub.conn.SetWriteDeadline(time.Now().Add(writeWait))

            es := event.String()

            logger.Info(20, file, f, fmt.Sprintf("Subscriber from %s receive event: %s",
                    sub.addr, es), nil)

            b, err := proto.Marshal(event)
            if err != nil {
                logger.Info(50, file, f, fmt.Sprintf("Subscriber from %s failed to serialize event %s",
                    sub.addr, es), err)
            }

            if err = sub.conn.WriteMessage(websocket.BinaryMessage, b); err != nil {
                logger.Info(40, file, f, fmt.Sprintf("Failed to write payload to Subscriber from %s. Remaining failed write attempt: %d",
                    sub.addr, attempt - 1), err)
                attempt--
                if attempt <= 0 {
                    logger.Info(50, file, f, fmt.Sprintf("Failed write attempt limit reaches for Subscriber from %s. Forcefully close connection", 
                            sub.addr), nil)
                    return
                }
			}
		}
	}
}

