package sockets

import (
	"io"
	"log"
	"net"
	"strconv"
)

func NewSockets(port string)  {
	sockets_proxy(port)
}

func sockets_proxy(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Panic(err)
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleRequest(conn)
	}
}
func handleRequest(conn net.Conn) {
	if conn == nil {
		return
	}
	defer conn.Close()
	var b [1024]byte
	n, err := conn.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	if b[0] == 0x05 {
		conn.Write([]byte{0x05, 0x00})
		n, err = conn.Read(b[:])
		var host, port string
		switch b[3] {
		case 0x01: //IP V4
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		case 0x03: //domain
			host = string(b[5 : n-2]) //b[4] domain length
		case 0x04: //IP V6
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))
		connection, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			log.Println(err)
			return
		}
		defer connection.Close()
		conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		go io.Copy(connection, conn)
		io.Copy(conn, connection)
	}
}