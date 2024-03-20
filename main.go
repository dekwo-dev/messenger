package main

import (
    "dekwo.dev/messager/route"
    . "dekwo.dev/messager/worker"
)

func main() {
    n := NewNotifier()
	d := NewDispatcher(n)

    go n.Run()
	go d.Run()

    route.Setup(d)
    route.RunTLS()
}
