package main

import (
	"fmt"
	"log"
	"net"
)

func checkErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    log.Println("Socket Server Running...")
    server, err := net.Listen("tcp", "localhost:8000")
    checkErr(err)
    defer server.Close()
    
    log.Println("Listening on localhost:8000\nWaiting for subscriber...")
    for {
        connection, err := server.Accept()
        checkErr(err)
        go func() {
            fmt.Print(connection.LocalAddr())
            fmt.Println(connection.RemoteAddr())
            connection.Close()
        }()
    }
}
