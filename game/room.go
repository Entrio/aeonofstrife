package game

import (
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
)

const (
	TILE_TYPE_WALL = uint8(iota)
	TILE_TYPE_DIRT
	TILE_TYPE_PORTAL
	TILE_TYPE_AIR
)

type Vector2 struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Tile struct {
	Type       uint8   `json:"type"`
	IsPassable bool    `json:"is_passable"`
	Position   Vector2 `json:"position"`
}

type Room struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Width          int            `json:"width"`
	Height         int            `json:"height"`
	Tiles          [][]Tile       `json:"-"`
	Entry          RoomEntryPoint `json:"entry"`
	Exit           RoomExitPoint  `json:"exit"`
	IsStartingRoom bool           `json:"is_starting_room"`
	isActive       bool           `json:"is_active"`
}

type RoomEntryPoint struct {
	Password       *string `json:"password"`
	LocationInRoom Vector2 `json:"location_in_room"`
}

type RoomExitPoint struct {
	LocationInRoom Vector2     `json:"location_in_room"`
	Destinations   []uuid.UUID `json:"destinations"`
}

func (rep RoomEntryPoint) GetPassword() *string {
	return rep.Password
}

func NewRoom(width, height int) *Room {
	rid := uuid.New()
	tiles := make([][]Tile, width)
	for k := range tiles {
		tiles[k] = make([]Tile, height)
	}
	return &Room{
		ID:          rid,
		Name:        fmt.Sprintf("room_%s", rid.String()),
		Description: "A generic room",
		Width:       width,
		Height:      height,
		Tiles:       tiles,
		Entry: RoomEntryPoint{
			LocationInRoom: Vector2{1, 1},
		},
		Exit: RoomExitPoint{
			LocationInRoom: Vector2{4, 4},
			Destinations:   nil,
		},
		IsStartingRoom: false,
	}
}

func (room *Room) UpdateTile(x, y uint8, tile Tile) *Room {
	fmt.Println(
		fmt.Sprintf(
			"Updating tile for room %s: %d",
			room.ID, tile.Type,
		),
	)
	room.Tiles[x][y] = tile
	return room
}

// Get the current room as a byte slice
func (room *Room) GetRoomByteData() (buffer []byte) {

	// This has the same structure as the room data packet for the network protocol
	sizeBuff := make([]byte, 2)

	// Get the length of the room and put it in temp buffer and append it to the general buffer
	binary.LittleEndian.PutUint16(sizeBuff, uint16(len(room.ID.String())))
	buffer = append(buffer, sizeBuff...)
	// Append the room id string to the buffer
	buffer = append(buffer, []byte(room.ID.String())...)

	return
}
