package game

import (
	"fmt"
	"net"
	"time"
)

func (server *Server) AddConnection(conn net.Conn) *Connection {
	newConnection := &Connection{
		conn:          conn,
		timeConnected: time.Now(),
		player:        nil,
	}

	server.connectionsList = append(server.connectionsList, newConnection)
	go newConnection.listen()
	welcomePacket := NewPacket(MsgWelcome)
	welcomePacket.WriteString("Welcome to the super awesome server This is a server message!")
	sendMessageToConnection(newConnection, *welcomePacket)
	return newConnection
}

/**
Handle player disconnects
*/
func (server *Server) onClientConnectionClosed(connection *Connection, err error) {
	for i, conn := range server.connectionsList {
		if conn == connection {
			// bye bye, remove from the slice and reshuffle
			server.connectionsList[i] = server.connectionsList[len(server.connectionsList)-1]
			server.connectionsList[len(server.connectionsList)-1] = nil
			server.connectionsList = server.connectionsList[:len(server.connectionsList)-1]
			fmt.Println(fmt.Sprintf("Disconnect from from %s", connection.conn.RemoteAddr().String()))
			break
		}
	}
}
