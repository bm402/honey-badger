package main

import (
	"log"
	"net"
	"time"
)

// listens for tcp connections on the given port
func listen(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Error creating listener on port ", port, ": ", err.Error())
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Error accepting connection on port", port, ": ", err.Error())
		}
		conn.SetDeadline(time.Now().Add(15 * time.Minute))
		go handle(conn, port)
	}
}

// serves a false command prompt on the given tcp connection and reads input
func handle(conn net.Conn, port string) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		_, err := conn.Write([]byte("$ "))
		if err != nil {
			break
		}

		rawInputLen, err := conn.Read(buf)
		if err != nil {
			break
		}
		input := string(buf[:rawInputLen])
		writeInputsToRawLogsTable(conn, port, input)
	}
}
