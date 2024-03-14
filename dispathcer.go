package main

import (
	"fmt"
)

type Dispatcher struct {
	subs    map[*Subscriber]E
    workers map[*Worker]E
    done    chan *Worker
	sub     chan *Subscriber // HTTP server   -> dispatcher
	unsub   chan *Subscriber // a subscriber <-> dispatcher
}

func dispatcher() *Dispatcher {
	return &Dispatcher{
		subs:    make(map[*Subscriber]E),
        workers: make(map[*Worker]E),
        done:    make(chan *Worker), 
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
                sub.addr,
                len(d.subs),
            ), nil)

			e := &ViewCountEvent{
                "ViewCountEvent",
				NewSub,
				uint8(len(d.subs)),
			}

			for sub := range d.subs {
				sub.event <- e
			}
		case sub := <-d.unsub:
			if _, ok := d.subs[sub]; ok {
				delete(d.subs, sub)

                info(30, f, file, fmt.Sprintf(
                    "Dispatcher removed subscriber from %s. Number of subscribers: %d",
                    sub.addr,
                    len(d.subs),
                ), nil)

				e := &ViewCountEvent{
                    "ViewCountEvent",
					DelSub,
					uint8(len(d.subs)),
				}

				for sub := range d.subs {
					sub.event <- e
				}
			}
		}
	}
}
