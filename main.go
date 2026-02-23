package main

import (
	"log"
	"net"
)

func main() {
	s := newServer()
	go s.run()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("unable to start server: %s", err)
	}
	defer listener.Close()
	log.Printf("started server on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection: %s", err)
			continue
		}

		go s.newClient(conn)
	}
}
