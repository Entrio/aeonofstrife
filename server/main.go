package main

import (
	"fmt"
	"github.com/Entrio/aeonofstrife/game"
	"net"
)

var (
	s *game.Server
)

/**
VERSION pre-alpha 0.1
@author Alexander Titarenko <westal@gmail.com>
*/
func main() {
	server, err := game.GetServer()
	if err != nil {
		panic(err)
	}
	s = server

	port := fmt.Sprintf("0.0.0.0:%d", s.GetPort())
	tcpAddr, err := net.ResolveTCPAddr("tcp4", port)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	s.Start()
	fmt.Println(fmt.Sprintf("Server started on port %s", listener.Addr().String()))
	listenForConnections(listener)

}

func listenForConnections(listener *net.TCPListener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		newDescriptor(conn)
	}
}

func newDescriptor(connection net.Conn) {
	fmt.Println(fmt.Sprintf("New connection from %s", connection.RemoteAddr().String()))
	s.AddConnection(connection)
}
