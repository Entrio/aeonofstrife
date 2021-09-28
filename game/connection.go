package game

import (
	"fmt"
	"net"
	"time"
)

type (
	Connection struct {
		conn          net.Conn
		timeConnected time.Time
		player        *Player
		isEditor      bool
	}
)

// listen goroutine listens for incoming data
func (connection *Connection) listen() {
	recvBuf := make([]byte, 1024*1024)
	packet := &Packet{
		buffer:     make([]byte, 0),
		cursor:     0,
		Connection: connection,
	}

	for {
		cLen, err := connection.conn.Read(recvBuf)
		if err != nil {
			// client disconnected
			connection.conn.Close()
			ServerInstance.onClientConnectionClosed(connection, err)
			return
		}
		// Need now to copy the read bytes amount and process it
		//processMessage(connection, recvBuf, len)

		temp := make([]byte, cLen)
		temp = recvBuf[:cLen]
		packet.Reset(handleData(temp, packet))
	}

}

func sendMessageToConnection(connection *Connection, packet Packet) {

	bytes, err := connection.conn.Write(packet.GetBytes())
	if err == nil {
		fmt.Println(fmt.Sprintf("Written %d bytes to stream", bytes))
	} else {
		fmt.Println(err)
	}
}

func (connection *Connection) sendBytes(data []byte) {
	connection.conn.Write(data)
}
func (connection *Connection) sendString(data string) {
	connection.sendBytes([]byte(data))
}
