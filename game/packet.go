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
		Type:   MSG_NULL_IOTA,
		buffer: data,
		cursor: 0,
	}
}

func (packet *Packet) GetMessageType() PacketType {
	if packet.Type == MSG_NULL_IOTA {
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
Read a certain amount of bytes from the buffer.
*/
func (packet *Packet) ReadBytes(length uint32) []byte {
	if packet.Length() <= (packet.cursor + length) {
		// we got something
		data := packet.buffer[packet.cursor : packet.cursor+length]

		// move the cursor
		packet.cursor += length

		return data
	} else {
		// TODO: Panic or some sort of error
		panic(fmt.Errorf("could not read bytes from packet as we are outside of the index. Packet length: %d, reading %d", packet.Length(), length))
	}
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
Write a byte to teh buffer. This method does conversion and checking automatically.
*/
func (packet *Packet) WriteIntByte(data int) *Packet {
	if data > 255 {
		panic(fmt.Errorf("error while converting inot to byte. WriteIntByte only supports value up to 255"))
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
	for _, tile := range room.TileList {
		if tile.Type != TILE_TYPE_DIRT {
			tiles = append(tiles, tile)
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
