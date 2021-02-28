package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"socks5/socks5"
	"strconv"
)

var (
	port     = flag.Int("-port", 0, "the proxy port")
	host     = flag.String("-host", "127.0.0.1", "the address listen on")
	Username = flag.String("-token", "", "your username")
	Password = flag.String("-password", "", "your password")
)

func main() {
	flag.Parse()

	if *port == 0 {
		if envPort := os.Getenv("PORT"); envPort != "" {
			if i, err := strconv.Atoi(envPort); err == nil {
				*port = i
			}
		}
	}
	if *port == 0 {
		*port = 1080
	}

	addr := net.TCPAddr{
		IP:   net.ParseIP(*host),
		Port: *port,
	}
	listener, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Printf("Server listen on: %s:%d\n", *host, *port)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Client address: ", conn.RemoteAddr())
		go socks5.Handle(conn)
	}
}
