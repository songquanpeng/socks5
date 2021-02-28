package handler

import (
	"encoding/binary"
	"fmt"
	"net"
)

var (
	Version    = []byte("\x05")[0]
	BufferSize = 256
)

var (
	ConnectCommand      = []byte("\x01")[0]
	BindCommand         = []byte("\x02")[0]
	UDPAssociateCommand = []byte("\x03")[0]
)

var (
	IPv4Address = []byte("\x01")[0]
	DomainName  = []byte("\x03")[0]
	IPv6Address = []byte("\x04")[0]
)

var (
	SuccessReply                 = []byte("\x00")[0]
	GeneralFailureReply          = []byte("\x01")[0]
	ConnectionNotAllowedReply    = []byte("\x02")[0]
	NetworkUnreachableReply      = []byte("\x03")[0]
	HostUnreachableReply         = []byte("\x04")[0]
	ConnectionRefusedReply       = []byte("\x05")[0]
	TTLExpiredReply              = []byte("\x06")[0]
	CommandNotSupportedReply     = []byte("\x07")[0]
	AddressTypeNotSupportedReply = []byte("\x08")[0]
)

var (
	Username = ""
	Password = ""
	NoAuth   = true
)

func SetUsernameAndPassword(username, password string) {
	Username = username
	Password = password
	NoAuth = false
}

func Handle(conn *net.TCPConn) {
	if authentication(conn) {
		processRequest(conn)
	}
}

func processRequest(conn *net.TCPConn) {
	/*
		Server replies to the client:
		+----+-----+-------+------+----------+----------+
		|VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
		+----+-----+-------+------+----------+----------+
		| 1  |  1  | X'00' |  1   | Variable |    2     |
		+----+-----+-------+------+----------+----------+
	*/
	dst := processConnectRequest(conn)
	if dst == "" {
		fmt.Println("Destination is blank!")
		return
	}
	replyType := SuccessReply
	remoteConn, err := net.Dial("tcp", dst)
	if err != nil {
		replyType = GeneralFailureReply
		fmt.Println(err)
		return
	}
	reservedField := []byte("\x00")[0]
	// TODO: use the real info
	addressType := IPv4Address
	address := make([]byte, 4)
	port := make([]byte, 2)

	// Build reply
	var reply []byte
	reply = append(reply, Version)
	reply = append(reply, replyType)
	reply = append(reply, reservedField)
	reply = append(reply, addressType)
	reply = append(reply, address...)
	reply = append(reply, port...)
	if _, err := conn.Write(reply); err != nil {
		fmt.Println(err)
		return
	}

	if replyType == SuccessReply {
		go forward(net.Conn(conn), remoteConn)
		go forward(remoteConn, net.Conn(conn))
	}
}

func processConnectRequest(conn *net.TCPConn) (destination string) {
	/*
		Client send its desired destination:
		+----+-----+-------+------+----------+----------+
		|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
		+----+-----+-------+------+----------+----------+
		| 1  |  1  | X'00' |  1   | Variable |    2     |
		+----+-----+-------+------+----------+----------+
	*/
	buffer := make([]byte, BufferSize)
	if _, err := conn.Read(buffer); err != nil {
		fmt.Println(err)
		return
	}
	if buffer[0] != Version {
		return
	}
	addressLength := 4
	address := ""
	switch buffer[1] {
	case ConnectCommand:
		switch buffer[3] {
		case IPv4Address:
			address = fmt.Sprintf("%d.%d.%d.%d", buffer[4], buffer[5], buffer[6], buffer[7])
		case IPv6Address:
			addressLength = 16
			address = parseIPv6Address(buffer[4 : 4+addressLength])
			address = "[" + address + "]"
		case DomainName:
			addressLength = int(buffer[4])
			address = string(buffer[5 : 5+addressLength])
			addressLength += 1 // Now it's equal to the length of DST.ADDR
		}
	case BindCommand:
		// TODO: support Bind command
	case UDPAssociateCommand:
		// TODO: support UDP Associate command
	}
	port := binary.BigEndian.Uint16(buffer[4+addressLength : 4+addressLength+2])
	return fmt.Sprintf("%s:%d", address, port)
}

func parseIPv6Address(bytes []byte) (addr string) {
	for i := 0; i < 16; i++ {
		if i != 0 && i%2 == 0 {
			addr += ":"
		}
		addr += fmt.Sprintf("%x", bytes[i])
	}
	return
}
