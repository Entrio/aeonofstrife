package game

import (
	"fmt"
	"net"
	"time"
)

type Connection struct {
	conn          net.Conn
	timeConnected time.Time
	player        *Player
	isEditor      bool
}

/**
Listen for incoming data
*/
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

func handleData(data []byte, receivedData *Packet) bool {
	packetLength := uint32(0)
	receivedData.SetBytes(data)

	if receivedData.UnreadLength() >= 4 {
		packetLength = receivedData.ReadUInt32()
		if packetLength <= 0 {
			return true
		}
	}

	for packetLength > 0 && packetLength <= receivedData.UnreadLength() {
		packetBytes := receivedData.ReadBytes(packetLength)

		//TODO: Maybe look into goroutine
		newPaket := NewUnknownPacket(packetBytes)
		newPaket.Connection = receivedData.Connection
		fmt.Println(fmt.Sprintf("Total handlers: %d", len(ServerInstance.packetHandler)))

		t := newPaket.GetMessageType()
		_, ok := ServerInstance.packetHandler[t]

		if !ok {
			fmt.Println(fmt.Sprintf("There is no handler registered for packet type %d", t))
		}

		ServerInstance.packetHandler[t].handle(newPaket)

		packetLength = 0
		if receivedData.UnreadLength() >= 4 {
			packetLength = receivedData.UnreadLength()
			if packetLength <= 0 {
				return true
			}
		}
	}

	if packetLength <= 1 {
		return true
	}

	return false
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
