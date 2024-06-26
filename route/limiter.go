package route

import (
	"fmt"
	"net/http"
	"strings"

	. "dekwo.dev/messager/env"
	"dekwo.dev/messager/logger"
	"golang.org/x/time/rate"

	"github.com/gorilla/websocket"
)

func checkOrigin(r *http.Request) bool {
    const f = "checkOrigin"
    const file = "route/limiter.go"

	o := r.Header.Get("Origin")

    logger.Info(20, file, f, fmt.Sprintf("Incoming origin: %s", o), nil)
    
	local := strings.HasPrefix(o, "http://localhost") ||
		strings.HasPrefix(o, "http://127.0.0.1")
    remote := strings.HasPrefix(
            o, fmt.Sprintf("https://www.%s/", GetEnv("FRONTEND_DOMAIN")),
        ) || 
              strings.HasPrefix(
            o, fmt.Sprintf("https://%s/", GetEnv("FRONTEND_DOMAIN")),
        )
    if Prod() {
        return remote
    } else {
        return local
    }
}

var upgrader = websocket.Upgrader {
	ReadBufferSize:  8,
	WriteBufferSize: 32,
	CheckOrigin: checkOrigin,
}

func rateLimited(rl rate.Limit, b int, next http.HandlerFunc) http.HandlerFunc {
    const f = "rateLimited"
    const file = "limiter.go"
    limiter := rate.NewLimiter(rl, b)

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if limiter.Allow() == false {
            logger.Info(40, file, f, "Too many WebSocket Request Connection", nil)
            http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
        }

        next.ServeHTTP(w, r)
    })
}
