package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	port := flag.String("p", "8081", "Port to listen on")
	flag.Parse()

	listen(*port)
}

func listen(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Error creating listener on port", port)
	}
	fmt.Println("Listening on port", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Error accepting connection on port", port)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	log.Println("Handling connection from", conn.RemoteAddr().String())
}
