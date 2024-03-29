package game

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Packet struct {
	Type       PacketType
	Connection *Connection
	buffer     []byte
	cursor     uint32
}

func NewPacket(packetType PacketType) *Packet {
	return &Packet{
		Type:   packetType,
		buffer: make([]byte, 0),
		cursor: 0,
	}
}

/**
Create a new packet with just bytes
*/
func NewUnknownPacket(data []byte) *Packet {

	return &Packet{
		Type:   MsgNullIota,
		buffer: data,
		cursor: 0,
	}
}

func (packet *Packet) GetMessageType() PacketType {
	if packet.Type == MsgNullIota {
		pt := PacketType(packet.ReadUint16())

		return pt
	}
	return packet.Type
}

// Get the packet as a byte array, ready for sending
func (packet *Packet) GetBytes() []byte {
	pSize := make([]byte, 4)
	pType := make([]byte, 2)
	binary.LittleEndian.PutUint16(pType, uint16(packet.Type))
	binary.LittleEndian.PutUint32(pSize, uint32(len(packet.buffer)+2)) // +2 for packet type

	pSize = append(pSize, pType...)
	return append(pSize, packet.buffer...)
}

// Read UUID 16 byte long string
func (packet *Packet) ReadUUID() string {
	length := packet.ReadUint16()
	if length != 36 {
		panic(fmt.Sprintf("Wrong UUID length! Expecting 36, got %d", length))
	}
	return string(packet.ReadBytes(uint32(length)))
}

/*
Set bytes for current packet. This is used when packets are fragmented.
*/
func (packet *Packet) SetBytes(data []byte) {
	packet.buffer = data
}

/*
Get the buffer size of the packet
*/
func (packet *Packet) Length() uint32 {
	return uint32(len(packet.buffer))
}

/**
Get the number of unread bytes in the packet buffer
*/
func (packet *Packet) UnreadLength() uint32 {
	if packet.cursor >= packet.Length() {
		return 0
	}
	return packet.Length() - packet.cursor
}

/**
Read a 32 bit unsigned integer. This is 4 bytes
*/
func (packet *Packet) ReadUInt32() (val uint32) {
	buf := bytes.NewBuffer(packet.buffer[packet.cursor : packet.cursor+4])
	if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
		panic(err)
		return
	}

	packet.cursor += 4
	return val
}

/**
Read a 16 bit unsigned integer. This is 2 bytes
*/
func (packet *Packet) ReadUint16() uint16 {
	val := uint16(0)
	val = binary.LittleEndian.Uint16(packet.buffer[packet.cursor : packet.cursor+2])
	packet.cursor += 2
	return val
}

/**
Read a 16 bit unsigned integer and return it as a signed 32 (or 64 depending on the system) bit integer
*/
func (packet *Packet) ReadUint16AsInt() int {
	return int(packet.ReadUint16())
}

/**
Read an 8 bit unsigned integer. This is 1 bytes
*/
func (packet *Packet) ReadUint8() uint8 {
	val := byte(0)
	val = packet.buffer[packet.cursor : packet.cursor+1][0]
	packet.cursor += 1
	return val
}

/**
Read a certain amount of bytes from the buffer.
*/
func (packet *Packet) ReadBytes(lenToRead uint32) []byte {

	// Check to see if we are not reading past the buffer lenToRead

	if packet.Length() >= (packet.cursor + lenToRead) {
		// we got something
		data := packet.buffer[packet.cursor : packet.cursor+lenToRead]

		// move the cursor
		packet.cursor += lenToRead

		return data
	} else {
		// TODO: Panic or some sort of error
		panic(
			fmt.Errorf(
				"Attempted to read outisde of slice bounds. Packet length: %d, reading %d bytes. Read index from %d - %d\nUnread: %d",
				packet.Length(), lenToRead, packet.cursor, packet.cursor+lenToRead, packet.UnreadLength(),
			),
		)
	}
}

/**
Read a single byte from the packet stream
*/
func (packet *Packet) ReadByte() byte {
	return packet.ReadBytes(1)[0]
}

/**
Read boolean value
*/
func (packet *Packet) ReadBoolean() bool {
	val := packet.ReadBytes(1)[0]
	if val == 1 {
		return true
	}
	return false
}

/**
Write a 36 character UUID string to the buffer.
*/
func (packet *Packet) WriteUUIDString(data string) {
	packet.WriteString(data)
}

/**
Write a string of undetermined length to the buffer. This writes the length first as an unsigned 16 bit integer
*/
func (packet *Packet) WriteString(data string) *Packet {
	if len(packet.buffer) == 0 {
		// nil slice,
		packet.buffer = make([]byte, 0)
	}

	packet.WriteUint16(uint16(len(data)))

	//fmt.Println(fmt.Sprintf("Written string: %s, %d bytes in length", data, uint16(len(data))))
	packet.buffer = append(packet.buffer, []byte(data)...)
	return packet
}

/**
Write a single byte to the buffer
*/
func (packet *Packet) WriteByte(data byte) {
	packet.buffer = append(packet.buffer, data)
}

/**
Write a byte to the buffer. This method does conversion and checking automatically.
Max value: 255
*/
func (packet *Packet) WriteIntByte(data int) *Packet {
	if data > 255 {
		panic(fmt.Errorf("error while converting int to byte. WriteIntByte only supports value up to 255"))
	}

	packet.WriteByte(byte(data))
	return packet
}

/**
This is an alias method for WriteByte
*/
func (packet *Packet) WriteUint8(data uint8) *Packet {
	packet.WriteByte(data)
	return packet
}

/**
Write a 2 byte, 16 bit unsigned integer to the buffer and offset the cursor.
*/
func (packet *Packet) WriteUint16(data uint16) *Packet {
	tBuffer := make([]byte, 2)

	binary.LittleEndian.PutUint16(tBuffer, data)
	packet.buffer = append(packet.buffer, tBuffer...)
	return packet
}

/**
Write a single byte as a boolean
*/
func (packet *Packet) WriteBool(data bool) *Packet {

	if data {
		packet.WriteByte(1)
	} else {
		packet.WriteByte(0)
	}

	return packet
}

/**
Parse the room into bits and write correct packets
*/
func (packet *Packet) WriteRoomData(room *Room) {
	packet.WriteString(room.ID.String())
	packet.WriteString(room.Name)
	packet.WriteString(room.Description)
	packet.WriteIntByte(room.Width)
	packet.WriteIntByte(room.Height)

	tiles := make([]Tile, 0)
	for x, b := range room.Tiles {
		for y := range b {
			tiles = append(tiles, room.Tiles[x][y])
		}
	}

	packet.WriteUint16(uint16(len(tiles)))

	for _, tile := range tiles {
		packet.WriteUint8(tile.Type).
			WriteBool(tile.IsPassable).
			WriteIntByte(tile.Position.X).
			WriteIntByte(tile.Position.Y)
	}

	packet.WriteUint8(uint8(room.Exit.LocationInRoom.X))
	packet.WriteUint8(uint8(room.Exit.LocationInRoom.Y))
	packet.WriteUint8(uint8(room.Entry.LocationInRoom.X))
	packet.WriteUint8(uint8(room.Entry.LocationInRoom.Y))
}

func (packet *Packet) Reset(force bool) {
	if force {
		packet.buffer = make([]byte, 0)
		packet.cursor = 0
	} else {
		packet.cursor -= 4
	}
}

func (packet *Packet) ResetCursor() *Packet {
	packet.cursor = 0
	return packet
}
