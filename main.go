/* multiple port version of port forwarding */
package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	conns := make(chan net.Conn)

	for port := 5000; port < 6000; port++ {
		listen(uint(port), conns)
	}

	for conn := range conns {
		log.Println("new client connected to", conn.RemoteAddr(), conn.LocalAddr())
		go handleRequest(conn)
	}
}

func listen(port uint, ch chan net.Conn) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("Error listening at port %d\n%v", port, err)
		return
	}
	log.Printf("listen at %d", port)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("Error accepting at port %v\n%v", port, err)
				return
			}
			ch <- conn
		}
	}()
}

func handleRequest(conn net.Conn) {
	fmt.Println("served at port", conn.LocalAddr())
	handlers := getHandlers()
	for _, handler := range handlers {
		handler.Do(conn)
	}
}

func getHandlers() []Handler {
	var h Upstream
	return []Handler{&h}
}

type Handler interface {
	Do(conn net.Conn)
}

type Upstream struct {
}

func (uh *Upstream) Do(conn net.Conn) {
	const upstream = "127.0.0.1:80"

	proxy, err := net.Dial("tcp", upstream)
	if err != nil {
		panic(err)
	}

	go copyIO(conn, proxy)
	go copyIO(proxy, conn)
}

func copyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}
