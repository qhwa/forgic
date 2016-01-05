/* multiple port version of port forwarding */
package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	conns := make(chan net.Conn)

	for port := 5000; port < 8000; port++ {
		go listen(port, conns)
	}

	for conn := range conns {
		fmt.Println("new client connected to", conn.RemoteAddr(), conn.LocalAddr())
		go handleRequest(conn)
	}
}

func listen(port uint, ch chan net.Conn) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		ch <- conn
	}
}

func handleRequest(conn net.Conn) {
	fmt.Println("served at port", conn.LocalAddr())

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
