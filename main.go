package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"socks5/handler"
	"strconv"
)

var (
	port     = flag.Int("port", 1080, "the proxy port")
	host     = flag.String("host", "localhost", "the address listen on")
	username = flag.String("username", "", "your username")
	password = flag.String("password", "", "your password")
)

func main() {
	flag.Parse()

	if *username != "" && *password != "" {
		handler.SetUsernameAndPassword(*username, *password)
	}

	if *port == 1080 {
		if envPort := os.Getenv("PORT"); envPort != "" {
			if i, err := strconv.Atoi(envPort); err == nil {
				*port = i
			}
		}
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
		go handler.Handle(conn)
	}
}
