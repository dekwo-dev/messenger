package route

import (
	"fmt"
	"log"
	"net/http"

    "dekwo.dev/messager/env"
    "dekwo.dev/messager/logger"
    "dekwo.dev/messager/worker"
)

func Setup(d *worker.Dispatcher) {
	http.HandleFunc("/ws", func(w http.ResponseWriter,
		r *http.Request) {
		ws(d, w, r)
	})
}

func Run() {
    const f = "Run"
    const file = "route/setup.go"

    logger.Info(20, file, f, fmt.Sprintf("Serving at port %d", env.Port()), nil)

    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", env.Port()), nil))
}
