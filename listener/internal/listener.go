package listener

import (
	"log"
	"net"
)

func listen(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Error creating listener on port", port)
	}
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
