package main

import (
	"flag"
	"fmt"
)

func main() {
	port := flag.String("p", "8081", "Port to listen on")
	fmt.Println("TODO: Honey Badger listener on port", *port)
}
