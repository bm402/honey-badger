package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	port := flag.String("p", "8081", "Port to listen on")
	flag.Parse()

	listen(*port)
}

// listens for tcp connections on the given port
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
		inputs := strings.Split(string(buf[:rawInputLen]), "\n")
		writeInputsToLog(conn, port, inputs)
	}
}

// writes inputs from the tcp connection to a log file
func writeInputsToLog(conn net.Conn, port string, inputs []string) {
	file, err := os.OpenFile("honey-badger-port-"+port+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)
	for _, input := range inputs {
		if len(input) > 0 {
			logger.Println("|", conn.RemoteAddr().String(), "|", input)
		}
	}
}
