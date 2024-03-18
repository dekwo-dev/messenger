package main

import (
	"fmt"

	"dekwo.dev/messager/pb"
)

type Dispatcher struct {
	subs    map[*Subscriber]E
    done    chan *Worker
    enqueue chan int
	sub     chan *Subscriber // HTTP server   -> dispatcher
	unsub   chan *Subscriber // a subscriber <-> dispatcher
}

func dispatcher(n *Notifier) *Dispatcher {
	return &Dispatcher{
		subs:    make(map[*Subscriber]E),
        done:    make(chan *Worker), 
        enqueue: n.enqueue,
		sub:     make(chan *Subscriber),
		unsub:   make(chan *Subscriber),
	}
}

func (d *Dispatcher) run() {
	const f = "Dispatcher.run"
    const file = "dispatcher.go"

	for {
		select {
		case sub := <-d.sub:
			d.subs[sub] = E{}

            info(30, f, file, fmt.Sprintf(
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

                info(30, f, file, fmt.Sprintf(
                    "Dispatcher removed subscriber from %s. Number of subscribers: %d",
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

				for sub := range d.subs {
					sub.event <- wrapper 
				}

                sub = nil
			}
        case load := <-d.enqueue:
            // TODO: Performance test between multiple workers and single worker 
            info(30, f, file, 
                fmt.Sprintf("Notifier notify %d items in queue", load), nil)

            w := worker(d)

            go w.run(load)
        case w := <- d.done:
            // TODO: Performance test on forwarding result between dispatcher and workers
            info(30, f, file,
                fmt.Sprintf("Received worker payload with a load of %d", 
                    len(w.payload)), nil)

            for _, event := range w.payload {
                info(30, f, file, fmt.Sprintf("Event: %v", event.String()), nil)
            }

            w = nil
		}
	}
}
