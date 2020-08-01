package game

import (
	"fmt"
	"github.com/google/uuid"
)

const (
	TILE_TYPE_WALL = uint8(iota)
	TILE_TYPE_DIRT
	TILE_TYPE_PORTAL
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
	TileList       []Tile         `json:"tile_list"`
	Entry          RoomEntryPoint `json:"entry"`
	Exit           RoomExitPoint  `json:"exit"`
	IsStartingRoom bool           `json:"is_starting_room"`
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
	for k, _ := range tiles {
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
