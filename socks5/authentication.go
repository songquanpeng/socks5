package socks5

import (
	"fmt"
	"net"
)

var (
	NoAuthenticationRequired = []byte("\x00")[0]
	GSSAPI                   = []byte("\x01")[0]
	UsernamePassword         = []byte("\x02")[0]
	// X'03' to X'7F' IANA ASSIGNED
	// X'80' to X'FE' RESERVED FOR PRIVATE METHODS
	NoAcceptableMethods = []byte("\xFF")[0]
)

func authentication(conn *net.TCPConn) (ok bool) {
	/*
		Server replies to the client:
		+----+--------+
		|VER | METHOD |
		+----+--------+
		| 1  |   1    |
		+----+--------+
	*/
	method := processMethodSelectionRequest(conn)
	reply := make([]byte, 2)
	reply[0] = Version
	reply[1] = method
	if _, err := conn.Write(reply); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func processMethodSelectionRequest(conn *net.TCPConn) (method byte) {
	/*
		Client send method selection request:
		+----+----------+----------+
		|VER | NMETHODS | METHODS  |
		+----+----------+----------+
		| 1  |    1     | 1 to 255 |
		+----+----------+----------+
	*/
	method = NoAcceptableMethods
	buffer := make([]byte, BufferSize)
	if _, err := conn.Read(buffer); err != nil {
		fmt.Println(err)
		return
	}
	// Check version
	if buffer[0] != Version {
		return NoAcceptableMethods
	}
	methodNum := int(buffer[1])
	methods := buffer[2 : 2+methodNum]
	if len(methods) != methodNum {
		return NoAcceptableMethods
	}
	for _, method := range methods {
		switch method {
		// TODO: support more methods
		case NoAuthenticationRequired:
			return method
			//case GSSAPI:
			//	return method
			//case UsernamePassword:
			//	return UsernamePassword
		}
	}
	return NoAcceptableMethods
}
