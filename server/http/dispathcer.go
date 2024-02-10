package main

import "log"

type Dispatcher struct {
    subs map[*Subscriber]struct{}
    event chan []byte // http server <-> dispathcer
    sub chan *Subscriber // http server <-> dispatcher
    unsub chan *Subscriber // a subscriber <-> dispatcher
}

func newDispatcher() *Dispatcher {
    return &Dispatcher{
        subs: make(map[*Subscriber]struct{}),
        event: make(chan []byte),
        sub: make(chan *Subscriber),
        unsub: make(chan *Subscriber),
    }
}

func (d *Dispatcher) run() {
    const f = "Dispatcher.run"
    log.Printf("%v (INFO): Dispatcher is running\n", f)
    for {
        select {
        case sub := <- d.sub:
            d.subs[sub] = struct{}{}
            log.Printf("%v (INFO): Dispatcher added subscriber from %v\n", 
                f, sub.addr)
            for sub := range d.subs {
                sub.event <- []byte(NewSub) 
            }
        case sub := <- d.unsub:
            if _, ok := d.subs[sub]; ok {
                delete(d.subs, sub)
                log.Printf("%v (INFO): Dispatcher removed subscriber from %v\n", 
                    f, sub.addr)
                for sub := range d.subs {
                    sub.event <- []byte(SubDisconnect)
                }
            }
        case event := <-d.event:
            log.Printf("%v (INFO): Dispatcher received event - %v from DB\n", 
                    f, string(event))
            switch string(event) {
            case NewComment, DelComment:
                for sub := range d.subs {
                    sub.event <- event
                }
            default:
                log.Printf("%v (WARN): Unknown event - %v from DB\n", 
                    f, string(event))
            }
        }
    }
}
