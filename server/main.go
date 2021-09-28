package main

import (
	"fmt"
	"net"
	"os"

	"github.com/Entrio/aeonofstrife/game"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	s *game.Server
)

/**
@version pre-alpha 0.1
@author Alexander Titarenko <westal@gmail.com>
*/
func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	var err error

	s, err = game.GetServer()
	if err != nil {
		panic(err)
	}

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
