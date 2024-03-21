package route

import (
	"fmt"
	"log"
	"net/http"

	"dekwo.dev/messager/env"
	"dekwo.dev/messager/logger"
	"dekwo.dev/messager/worker"
	"golang.org/x/crypto/acme/autocert"
)

func Setup(d *worker.Dispatcher) {
    http.Handle("/ws", rateLimited(2, 4,
        func (w http.ResponseWriter, r *http.Request) {
            ws(d, w, r)
        },
    ))
}

func Run() {
    const f = "Run"
    const file = "route/setup.go"

    logger.Info(20, file, f, fmt.Sprintf("Serving at port %d", env.Port()), nil)

    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", env.Port()), nil))
}

func RunTLS() {
    const f = "RunTLS"
    const file = "route/setup.go"

    domain := "localhost"
    if env.Prod() {
        domain = "dekwo.dev"
    }

    m := &autocert.Manager {
        Cache: autocert.DirCache(".autocert"),
        Prompt: autocert.AcceptTOS,
        Email: "dekr0.dk@protonmail.com",
        HostPolicy: autocert.HostWhitelist(domain),
    }

    s := &http.Server {
        Addr: ":https", // Port 443
        TLSConfig: m.TLSConfig(),
    }

    logger.Info(20, file, f, "Serving at port HTTPS", nil)

    log.Fatal(s.ListenAndServeTLS("", ""))
}
