package main

import (
	"fmt"
	"log"
	"net"

	"example.com/event-emitter/config"
	"example.com/event-emitter/utils"
)

func main() {
    log.Println("Socket Server Running...")
    ip := utils.GetSelfPublicIP()
    server, err := net.Listen("tcp", ip + config.Port)
    utils.CheckError(err)
    defer server.Close()
    
    log.Println("Listening on " + ip + config.Port + " \nWaiting for subscriber...")
    for {
        connection, err := server.Accept()
        utils.CheckError(err)
        go func() {
            fmt.Print(connection.LocalAddr())
            fmt.Println(connection.RemoteAddr())
            connection.Close()
        }()
    }
}
