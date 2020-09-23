package tools

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Entrio/aeonofstrife/game"
	"math"
)

type ByteBuffer struct {
	buffer []byte
	cursor int
}

func NewByteBuffer() *ByteBuffer {
	return &ByteBuffer{}
}

/**
Get the formed buffer as a byte slice
*/
func (bf *ByteBuffer) GetBytes() (data []byte, length int) {
	data = make([]byte, len(bf.buffer))

	length = copy(data, bf.buffer)

	return
}

/**
Add given byte slice to the existing buffer
*/
func (bf *ByteBuffer) AddBytes(data []byte) *ByteBuffer {

	bf.buffer = append(bf.buffer, data...)

	return bf
}

/**
Add a string to the buffer. The second parameter indicates whether we should add 2 byte uint16
before the string to indicate the string length.
*/
func (bf *ByteBuffer) AddString(data string, allocSpace bool) *ByteBuffer {
	var strLength uint16

	if len(data) > math.MaxUint16 {
		panic(fmt.Errorf("string length exceeds max allowed length of %d characters", math.MaxUint16))
	}

	strLength = uint16(len(data))

	if allocSpace {
		bf.AddUint16(strLength)
	}

	bf.buffer = append(bf.buffer, data...)

	return bf
}

/**
Add a single byte to the slice
*/
func (bf *ByteBuffer) AddByte(data byte) *ByteBuffer {
	bf.buffer = append(bf.buffer, data)

	return bf
}

/**
Add a single boolean (byte) to the slice
*/
func (bf *ByteBuffer) AddBool(data bool) *ByteBuffer {

	if data {
		bf.buffer = append(bf.buffer, 1)
	} else {
		bf.buffer = append(bf.buffer, 0)
	}

	return bf
}

/**
Add an uint8 to the buffer array. An alias of AddByte
*/
func (bf *ByteBuffer) AddUint8(data uint8) *ByteBuffer {
	bf.AddByte(data)

	return bf
}

/**
Add an uint16 to the buffer array.
*/
func (bf *ByteBuffer) AddUint16(data uint16) *ByteBuffer {
	tBuffer := make([]byte, 2)

	binary.LittleEndian.PutUint16(tBuffer, data)
	bf.buffer = append(bf.buffer, tBuffer...)

	return bf
}

/**
Add an uint32 to the buffer array.
*/
func (bf *ByteBuffer) AddUint32(data uint32) *ByteBuffer {
	buff := new(bytes.Buffer)

	if err := binary.Write(buff, binary.LittleEndian, data); err != nil {
		panic(err)
	}

	bf.buffer = append(bf.buffer, buff.Bytes()...)

	return bf
}

/**
Add an uint64 to the buffer array.
*/
func (bf *ByteBuffer) AddUint64(data uint64) *ByteBuffer {
	buff := new(bytes.Buffer)

	if err := binary.Write(buff, binary.LittleEndian, data); err != nil {
		panic(err)
	}

	bf.buffer = append(bf.buffer, buff.Bytes()...)

	return bf
}

/**
Add an int8 to the buffer array.
*/
func (bf *ByteBuffer) AddInt8(data int8) *ByteBuffer {
	buff := new(bytes.Buffer)

	if err := binary.Write(buff, binary.LittleEndian, data); err != nil {
		panic(err)
	}

	bf.buffer = append(bf.buffer, buff.Bytes()...)

	return bf
}

/**
Add an int16 to the buffer array.
*/
func (bf *ByteBuffer) AddInt16(data int16) *ByteBuffer {
	buff := new(bytes.Buffer)

	if err := binary.Write(buff, binary.LittleEndian, data); err != nil {
		panic(err)
	}

	bf.buffer = append(bf.buffer, buff.Bytes()...)

	return bf
}

/**
Add an int32 to the buffer array.
*/
func (bf *ByteBuffer) AddInt32(data int32) *ByteBuffer {
	buff := new(bytes.Buffer)

	if err := binary.Write(buff, binary.LittleEndian, data); err != nil {
		panic(err)
	}

	bf.buffer = append(bf.buffer, buff.Bytes()...)

	return bf
}

/**
Add an int64 to the buffer array.
*/
func (bf *ByteBuffer) AddInt64(data int64) *ByteBuffer {
	buff := new(bytes.Buffer)

	if err := binary.Write(buff, binary.LittleEndian, data); err != nil {
		panic(err)
	}

	bf.buffer = append(bf.buffer, buff.Bytes()...)

	return bf
}

func (bf *ByteBuffer) GetRoomBytes(room *game.Room) ([]byte, int) {
	bf.AddString(room.ID.String(), true).
		AddString(room.Name, true).
		AddString(room.Description, true).
		AddByte(byte(room.Width)).
		AddByte(byte(room.Height))

	tiles := make([]game.Tile, 0)
	for x, b := range room.Tiles {
		for y := range b {
			tiles = append(tiles, room.Tiles[x][y])
		}
	}

	bf.AddUint16(uint16(len(tiles)))

	for _, tile := range tiles {
		bf.AddUint8(tile.Type).
			AddBool(tile.IsPassable).
			AddByte(byte(tile.Position.X)).
			AddByte(byte(tile.Position.Y))
	}

	bf.AddUint8(uint8(room.Exit.LocationInRoom.X))
	bf.AddUint8(uint8(room.Exit.LocationInRoom.Y))
	bf.AddUint8(uint8(room.Entry.LocationInRoom.X))
	bf.AddUint8(uint8(room.Entry.LocationInRoom.Y))

	return bf.GetBytes()
}
