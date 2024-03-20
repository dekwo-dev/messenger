package worker 

import (
	"fmt"

    . "dekwo.dev/messager/logger"
	"dekwo.dev/messager/pb"
)

type E struct {

}

type Dispatcher struct {
	subs    map[*Subscriber]E
    done    chan *Worker
    enqueue chan int
	Sub     chan *Subscriber // HTTP server   -> dispatcher
	unsub   chan *Subscriber // a subscriber <-> dispatcher
}

func NewDispatcher(n *Notifier) *Dispatcher {
	return &Dispatcher{
		subs:    make(map[*Subscriber]E),
        done:    make(chan *Worker), 
        enqueue: n.enqueue,
		Sub:     make(chan *Subscriber),
		unsub:   make(chan *Subscriber),
	}
}

func (d *Dispatcher) Run() {
	const f = "Dispatcher.Run"
    const file = "worker/dispatcher.go"

	for {
		select {
		case sub := <-d.Sub:
            if _, in := d.subs[sub]; in {
                Info(50, file, f, fmt.Sprintf("A subscriber from %s is reusing a connection. Dispatcher close the connection", 
                    sub.addr), nil)
                sub.conn.Close() // Behavior testing on read() and notify() thread
                return
            }
			d.subs[sub] = E{}

            Info(20, file, f, fmt.Sprintf(
                "Dispatcher added subscriber from %s. Number of subscribers: %d",
                sub.addr, len(d.subs)), nil)

            event := &pb.ViewCountEvent {
                ViewCountEventType: pb.ViewCountEventType_VIEWCOUNTEVENTTYPE_INCREASE,
                ViewCount: uint32(len(d.subs)),
            }

            wrapper := &pb.Event {
                EventOneof: &pb.Event_ViewCountEvent { ViewCountEvent: event }, 
            }

			for sub := range d.subs {
				sub.event <- wrapper
			}
		case sub := <-d.unsub:
			if _, ok := d.subs[sub]; ok {
				delete(d.subs, sub)

                Info(20, file, f, fmt.Sprintf("Dispatcher removed subscriber from %s. Number of subscribers: %d",
                    sub.addr, len(d.subs)), nil)

                event := &pb.ViewCountEvent {
                    ViewCountEventType: pb.ViewCountEventType_VIEWCOUNTEVENTTYPE_DECREASE,
                    ViewCount: uint32(len(d.subs)),
                }

                wrapper := &pb.Event {
                    EventOneof: &pb.Event_ViewCountEvent { 
                        ViewCountEvent: event,
                    },
                }

				for s := range d.subs {
					s.event <- wrapper 
				}

                sub = nil
			} else {
                Info(50, file, f, fmt.Sprintf("The connection of a subscriber from %s is being tracked by dispatcher",
                    sub.addr), nil)
            }
        case load := <-d.enqueue:
            // TODO: Performance test between multiple workers and single worker 
            Info(20, file, f, fmt.Sprintf("Notifier notify %d items in queue",
                load), nil)

            w := newWorker(d)

            go w.run(load)
        case w := <- d.done:
            // TODO: Performance test on forwarding result between dispatcher and workers
            Info(20, file, f, fmt.Sprintf("Received worker payload with a load of %d",
                len(w.payload)), nil)

            for _, event := range w.payload {
                for s := range d.subs {
                    s.event <- event
                }
            }
            w = nil
		}
	}
}
