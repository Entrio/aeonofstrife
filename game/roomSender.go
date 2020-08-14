package game

import (
	"encoding/binary"
	"fmt"
)

const (
	clientBufferSize = 325933 // 3 fully packed rooms
)

func sendRoomDataToConnection(connection *Connection) {
	fmt.Println(fmt.Sprintf("Sending rooms: %d", len(ServerInstance.roomList)))
	// return a msg of how many rooms

	/*
		for _, room := range ServerInstance.roomList {
			sendMessageToConnection(connection, MSG_ROOM_COUNT_RESPONSE, room.getRoomPacket())
		}
	*/

	/*
		for _, r := range ServerInstance.roomList {
			nr, err := json.Marshal(r)
			if err == nil {
				sendMessageToConnection(connection, MSG_ROOM_COUNT_RESPONSE, nr)
			}
		}
	*/
}

func buildRoomPacket() ([]byte, int) {
	builder := make([]byte, clientBufferSize)
	offset := 0
	for _, room := range ServerInstance.roomList {
		// 36 bytes room UUID
		rid := []byte(room.ID.String())
		for i := 0; i < 36; i++ {
			builder[i] = rid[i] // Optimize malloc
			offset++
		}

		// 1 byte room name
		builder[offset] = byte(len(room.Name))
		offset++

		// <n> bytes room name
		roomLen := len(room.Name)
		rd := []byte(room.Name)
		for i := 0; i < roomLen; i++ {
			builder[offset] = rd[i] // optimize malloc
			offset++
		}

		// 2 bytes room description
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(len(room.Description)))
		builder[offset] = b[0]
		offset++
		builder[offset] = b[1]
		offset++

		// <n> bytes room description
		roomDescLen := len(room.Description)
		roomDesc := []byte(room.Description)
		for i := 0; i < roomDescLen; i++ {
			builder[offset] = roomDesc[i] // optimize malloc
			offset++
		}

		// 1 byte room width
		builder[offset] = byte(room.Width)
		offset++

		// 1 byte room height
		builder[offset] = byte(room.Height)
		offset++

		tiles := make([]Tile, 0)
		for _, tile := range room.TileList {
			if tile.Type != TILE_TYPE_DIRT {
				tiles = append(tiles, tile)
			}
		}

		b = make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(len(tiles)))
		builder[offset] = b[0]
		offset++
		builder[offset] = b[1]
		offset++

		// Type
		// Passable
		// X
		// Y

		for _, tile := range tiles {
			builder[offset] = byte(tile.Type)
			offset++
			if tile.IsPassable {
				builder[offset] = 1
			} else {
				builder[offset] = 0
			}
			offset++
			builder[offset] = byte(tile.Position.X)
			offset++
			builder[offset] = byte(tile.Position.Y)
			offset++
		}

		fmt.Println(fmt.Sprintf("Offset stopped at %d", offset))
		fmt.Println(builder[:offset])
		fmt.Println(string(builder[:offset]))
	}
	return builder, offset - 1
}

//TODO Create a room packet function with automatic offset

/**
So the room packet should look something like this:
1 byte - number of rooms we are sending (max 255)

36 bytes - UUID
1 bytes - room name length uint8
<n> bytes - room string (max of 255 characters in UTF-8 encoding)
2 bytes - room description length uint16 (max 65535 characters in UTF-8 encoding)
<n> bytes - room description string (max 65535 characters in UTF-8 encoding)
1 byte - room width uint8 (max 255)
1 byte - room height uint8 (max 255)

At this stage it gets interesting, we assume that all tiles are walkable and are made of dirt except the ones we are sending
so if we take a 10x10 room and send 0 tiles, all tiles are walkable (and the game will break)

2 bytes - how many tiles we are sending uint16 (max 65535)
... 1 byte - tile type uint8 (max 255)
... 1 byte - passable uint8 (boolean)
... 1 byte - positionX uint8 (max 255)
... 1 byte - positionY uint8 (max 255)

1 byte - room exit position X
1 byte - room exit position Y
1 byte - room entry position X
1 byte - room entry position Y

and that's it, simple right. Maximum packet size: 325,933 bytes (can fit 3 rooms int 1 mb buffer)
16 + 1 + 255 + 2 + 65535 + 1 + 1 + 2 + (65025 * 4)

*/

func (room *Room) getRoomPacket() []byte {
	offset := 0
	ba := make([]byte, 0)
	roomID := []byte(room.ID.String())
	ba = append(ba, roomID...)
	offset += len(roomID)

	roomNameData := []byte(room.Name)
	var roomNameLength byte
	roomNameLength = uint8(len(room.Name))
	ba = append(ba, roomNameLength)
	offset++
	ba = append(ba, roomNameData...)
	offset += len(roomNameData)

	roomDescLength := make([]byte, 2)
	roomDescData := make([]byte, len(room.Description))
	binary.LittleEndian.PutUint16(roomDescLength, uint16(len(room.Description)))
	ba = append(ba, roomDescLength...)
	offset += 2
	ba = append(ba, roomDescData...)
	offset += len(room.Description)

	var roomWidth byte
	roomWidth = uint8(room.Width)
	ba = append(ba, roomWidth)
	offset++

	var roomHeight byte
	roomHeight = uint8(room.Height)
	ba = append(ba, roomHeight)
	offset++

	tiles := make([]Tile, 0)
	for _, tile := range room.TileList {
		if tile.Type != TILE_TYPE_DIRT {
			tiles = append(tiles, tile)
		}
	}

	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(len(tiles)))
	ba = append(ba, b...)
	offset += 2

	for _, tile := range tiles {
		ba = append(ba, byte(tile.Type))
		offset++
		if tile.IsPassable {
			ba = append(ba, 1)
		} else {
			ba = append(ba, 0)
		}
		offset++

		ba = append(ba, byte(tile.Position.X))
		offset++
		ba = append(ba, byte(tile.Position.Y))
		offset++
	}

	return ba
}
