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

func (connection *Connection) sendBytes(data []byte) {
	connection.conn.Write(data)
}
func (connection *Connection) sendString(data string) {
	connection.sendBytes([]byte(data))
}

/**
Listen for incoming data
*/
func (connection *Connection) listen() {
	recvBuf := make([]byte, 1024)
	for {
		len, err := connection.conn.Read(recvBuf)
		if err != nil {
			// client disconnected
			connection.conn.Close()
			ServerInstance.onClientConnectionClosed(connection, err)
			return
		}
		processMessage(connection, recvBuf, len)
	}
}

func processMessage(connection *Connection, data []byte, l int) {
	// Determine msg type
	msgType := binary.LittleEndian.Uint16(data[:2])

	switch msgType {
	case MSG_PING_RESPONSE:
		fmt.Println("GET A PING RESPONSE")
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

func sendMessageToConnection(connection *Connection, mType uint16, data []byte) {
	response := make([]byte, 2)
	size := make([]byte, 4)
	binary.LittleEndian.PutUint16(response, mType)
	binary.LittleEndian.PutUint32(size, uint32(len(data)))
	fmt.Println(fmt.Sprintf("Type: %d, size: %d", mType, uint32(len(data))))
	response = append(response, size...)
	response = append(response, data...)
	bytes, err := connection.conn.Write(response)
	if err == nil {
		fmt.Println(fmt.Sprintf("Written %d bytes to stream", bytes))
	} else {
		fmt.Println(err)
	}
}
