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

	port := fmt.Sprintf("%s:%d", s.GetAddress(), s.GetPort())
	tcpAddr, err := net.ResolveTCPAddr("tcp4", port)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to resolve TCP address")
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen on interface")
	}
	defer func() {
		err = listener.Close()
		if err != nil {
			log.Warn().Err(err).Msg("Failed to close listeners")
		}
	}()

	s.Start()
	log.Info().Str("address", listener.Addr().String()).Msg("Server started")
	listenForConnections(listener)

}

func listenForConnections(listener *net.TCPListener) {
	log.Debug().Msg("TCP Listener now is listening for connections")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Warn().Err(err).Msg("Failed to accept connection")
			continue
		}
		newDescriptor(conn)
	}
}

func newDescriptor(connection net.Conn) {
	log.Info().Str("connection", connection.RemoteAddr().String()).Msg("Incoming network connection")
	s.AddConnection(connection)
}
