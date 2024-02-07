package main

import (
	"io"
	"log"
	"net"
	"sync"

	"example.com/event-emitter/config"
	"example.com/event-emitter/utils"
)

type SafeSubsribers struct {
	mu          sync.Mutex
	subscribers map[string]net.Conn
}

const NEWCOMMENT = "on-new-comment"
const COMMENTDEL = "on-comment-delete"

// Should change this to connection pool implmentation
var subscribers = SafeSubsribers{subscribers: make(map[string]net.Conn)}

func getServerWatcher() net.Listener {
	log.Println("getServerMonitor: Launching HTTP server monitor...")
	watcher, err := net.Listen("tcp", "localhost:8888")
	utils.CheckError(err)
	log.Println(
		"getServerMonitor: HTTP server monitor is listening on localhost:8888",
	)
	return watcher
}

func getSubHandler() net.Listener {
	log.Println("getEventEmitter: Launching event emitter...")
	ip := utils.GetSelfPublicIP()
	emitter, err := net.Listen("tcp", ip+config.Port)
	utils.CheckError(err)

	log.Println(
		"getEventEmitter: Event emitter is listening on " + ip + config.Port,
	)
	log.Println("getEventEmitter: Waiting for upcoming subscriber")

	return emitter
}

func watchServer(watcher net.Listener) {
	defer watcher.Close()
    buffer := make([]byte, 32, 32)
	for {
		conn, err := watcher.Accept()
		utils.CheckError(err)
		addr := conn.RemoteAddr().String()
		log.Printf("watchServer: Connected to HTTP server at %v\n", addr)
		for { // watch until HTTP close down
			n, err := conn.Read(buffer)
			if err != nil {
				if err == io.EOF {
					log.Printf(
						"watchServer: Gracefully close connection with server at %v\n",
						addr,
					)
				} else {
					log.Printf("watchServer (WARN): %v\n", err)
					log.Printf("watchServer (WARN): Force closing connection with %v", addr)
				}
				conn.Close()
				break
			}

            event := string(buffer[:n])
			switch event {
			case NEWCOMMENT, COMMENTDEL:
				emitEvent(event)
			default:
				log.Printf("watchServer: Unknown event %v", event)
                log.Printf("watchServer: %v - %v", NEWCOMMENT, len(NEWCOMMENT),)
			}
            clear(buffer)
		}
	}
}

func emitEvent(event string) {
	subscribers.mu.Lock()

	test := make([]byte, 0, 1)
	for addr := range subscribers.subscribers {
		conn := subscribers.subscribers[addr]
		_, err := conn.Read(test)
		if err != nil {
			if err != io.EOF {
				log.Printf(
					"emitEvent (WARN): Error when checking %v connection\n",
					addr,
				)
				log.Printf(
					"emitEvent (WARN): Force closing connection with %v\n",
					addr,
				)
			} else {
				log.Printf("emitEvent: Gracefully close connection with %v\n",
					addr)
			}
			break
		}
		_, err = conn.Write([]byte(event))
		if err != nil {
			log.Printf(
				"emitEvent (WARN): Error when emitting event to %v\n",
				addr,
			)
		}
	}

	subscribers.mu.Unlock()
}

func main() {
	watcher := getServerWatcher()
	handler := getSubHandler()

	go watchServer(watcher)

	for {
		conn, err := handler.Accept()
		utils.CheckError(err)

		addr := conn.RemoteAddr().String()

		subscribers.mu.Lock()

		subscribers.subscribers[addr] = conn

		log.Printf("Accept one subscriber from address %v", addr)
		log.Printf("Number of subscribers: %d", len(subscribers.subscribers))

		subscribers.mu.Unlock()
	}
}
