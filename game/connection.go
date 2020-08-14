package game

import (
	"encoding/binary"
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
		ServerInstance.packetHandler[newPaket.GetMessageType()].handle(newPaket)

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
func processMessage(connection *Connection, data []byte, l int) {
	// Determine msg type
	// Construct a packet
	msgType := PacketType(binary.LittleEndian.Uint16(data[:2]))

	switch msgType {
	case MSG_PING_RESPONSE:
		fmt.Println("GOT A PING RESPONSE")
		break
	case MSG_ROOM_COUNT_REQUEST:
		sendRoomDataToConnection(connection)
	case MSG_ROOM_UPDATE_NAME:
		processRoomUpdateName(connection, data[2:l])
		break
	default:
		break
	}
}
func processRoomUpdateName(connection *Connection, data []byte) {
	roomID := string(data[:36])
	roomName := string(data[36:])

	r := ServerInstance.FindRoom(roomID)
	if r == nil {
		fmt.Println("No such room found!")
		return
	}

	r.Name = roomName

	fmt.Println(fmt.Sprintf("Updating room (%s) name to: %s", roomID, roomName))
}
