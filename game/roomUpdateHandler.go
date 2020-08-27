package game

import "fmt"

type RoomUpdateHandler struct{}

/*
*****************************
ROOM UPDATE PAYLOAD STRUCTURE
*****************************
1 byte - update type (0 - tiles)
16 bytes - string - room ID
2 bytes - uint16 number fo tiles that we will be sending
... 1 byte - tile type uint8 (max 255)
... 1 byte - is passable byte
... 1 byte - positionX uint8 (max 255)
... 1 byte - positionY uint8 (max 255)
*/

func (r RoomUpdateHandler) handle(packet *Packet) {
	updateType := packet.ReadBytes(1)
	roomID := packet.ReadUUID()
	tileCount := packet.ReadUint16AsInt()

	room := ServerInstance.FindRoom(roomID)
	if room == nil {
		panic(fmt.Sprintf("Failed to find room with UUID %s", roomID))
	}

	fmt.Println(
		fmt.Sprintf(
			"update type: %d for room %s. Tiles: %d",
			updateType, roomID, tileCount,
		),
	)

	for i := 0; i < tileCount; i++ {
		// pew pew lasers
		_tileType := uint8(packet.ReadByte())
		_isPassable := packet.ReadBoolean()
		_posX := packet.ReadUint8()
		_posY := packet.ReadUint8()

		fmt.Println(
			fmt.Sprintf(
				"Tile: %d, pass: %t X: %d, Y: %d",
				_tileType, _isPassable, _posX, _posY,
			),
		)
	}

}
