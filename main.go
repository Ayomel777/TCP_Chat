package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	s := newServer()
	go s.run()
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("enter the ip:\n")
	ip, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("error in reading IP")
		return
	}
	ip = strings.TrimSpace(ip)
	fmt.Printf("enter the port:\n")
	portStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("error in reading port\n")
		return
	}
	portStr = strings.TrimSpace(portStr)
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		fmt.Printf("incorrect port\n")
		return
	}
	addr := fmt.Sprintf("%s:%d", ip, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("unable to start server: %s", err)
	}
	defer listener.Close()
	log.Printf("started server on %s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection: %s", err)
			continue
		}

		go s.newClient(conn)
	}
}
