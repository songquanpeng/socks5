package handler

import (
	"fmt"
	"net"
)

var (
	NoAuthenticationRequired = []byte("\x00")[0]
	GSSAPI                   = []byte("\x01")[0]
	UsernamePassword         = []byte("\x02")[0]
	NoAcceptableMethods      = []byte("\xFF")[0]
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
	if method == NoAcceptableMethods {
		return false
	}
	reply := make([]byte, 2)
	reply[0] = Version
	reply[1] = method
	if _, err := conn.Write(reply); err != nil {
		fmt.Println(err)
		return false
	}
	if method == UsernamePassword {
		return usernamePasswordNegotiation(conn)
	}
	return true
}

func usernamePasswordNegotiation(conn *net.TCPConn) (ok bool) {
	/*
		Client send request:
		+----+------+----------+------+----------+
		|VER | ULEN |  UNAME   | PLEN |  PASSWD  |
		+----+------+----------+------+----------+
		| 1  |  1   | 1 to 255 |  1   | 1 to 255 |
		+----+------+----------+------+----------+
	*/
	ok = false
	version := uint8(1)
	buffer := make([]byte, 102)
	if _, err := conn.Read(buffer); err != nil {
		fmt.Println(err)
		return
	}
	if buffer[0] != version {
		return
	}
	usernameLength := buffer[1]
	username := string(buffer[2 : 2+usernameLength])
	passwordLength := buffer[2+usernameLength]
	password := string(buffer[2+usernameLength+1 : 2+usernameLength+1+passwordLength])

	/*
		Server reply:
		+----+--------+
		|VER | STATUS |
		+----+--------+
		| 1  |   1    |
		+----+--------+
	*/
	status := uint8(1)
	if username == Username && password == Password {
		status = 0
		ok = true
	}
	var reply []byte
	reply = append(reply, Version)
	reply = append(reply, status)
	if _, err := conn.Write(reply); err != nil {
		fmt.Println(err)
		return
	}
	return
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
		if NoAuth && method == NoAuthenticationRequired {
			return NoAuthenticationRequired
		}
		if !NoAuth && method == UsernamePassword {
			return UsernamePassword
		}
	}
	return NoAcceptableMethods
}
