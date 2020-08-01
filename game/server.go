package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"path"
	"time"
)

var ServerInstance *Server

type Server struct {
	connectionsList []*Connection
	roomList        map[string]*Room
	ticker          *time.Ticker
	config          *serverConfig
}

type serverConfig struct {
	ServerName string `json:"server_name"`
	ServerPort int    `json:"server_port"`
	RoomData   struct {
		Config struct {
			MinWidth  int `json:"min_width"`
			MaxWidth  int `json:"max_width"`
			MinHeight int `json:"min_height"`
			MaxHeight int `json:"maxHeight"`
		} `json:"config"`
		MinRooms int `json:"min_rooms"`
	} `json:"room_data"`
}

/**
Get an existing instance of create a new one
*/
func GetServer() (*Server, error) {
	if ServerInstance == nil {
		ServerInstance = &Server{
			connectionsList: make([]*Connection, 0),
			roomList:        map[string]*Room{},
		}
	}

	// Get current working directory
	cwd, _ := os.Getwd()
	fmt.Println(fmt.Sprintf("Setting working directory to: %s", cwd))

	dirs := checkDirectories(cwd)
	config, err := checkServerConfig(dirs[0])

	if err != nil {
		return nil, err
	}
	ServerInstance.config = config

	err = loadServerRooms(dirs[1])

	if err != nil {
		return nil, err
	}
	return ServerInstance, nil
}

func (server *Server) GetPort() int {
	return server.config.ServerPort
}

func (server *Server) Start() {
	server.ticker = time.NewTicker(time.Millisecond * 3000)
}

func (server *Server) GetName() string {
	return server.config.ServerName
}

/**
Find a room based on UUID string
*/
func (server *Server) FindRoom(uuid string) *Room {
	for i, room := range server.roomList {
		if room.ID.String() == uuid {
			return server.roomList[i]
		}
	}
	return nil
}

/**
Handle player disconnects
*/
func (server *Server) onClientConnectionClosed(connection *Connection, err error) {
	for i, conn := range server.connectionsList {
		if conn == connection {
			// bye bye, remove from the slice and reshuffle
			server.connectionsList[i] = server.connectionsList[len(server.connectionsList)-1]
			server.connectionsList[len(server.connectionsList)-1] = nil
			server.connectionsList = server.connectionsList[:len(server.connectionsList)-1]
			fmt.Println(fmt.Sprintf("Disconnect from from %s", connection.conn.RemoteAddr().String()))
			break
		}
	}
}

/**
New incoming connection
*/
func (server *Server) AddConnection(conn net.Conn) *Connection {
	newConnection := &Connection{
		conn:          conn,
		timeConnected: time.Now(),
		player:        nil,
	}

	server.connectionsList = append(server.connectionsList, newConnection)
	go newConnection.listen()
	sendMessageToConnection(newConnection, MSG_SPECIAL_2, []byte("Welcome to the server"))
	return newConnection
}

/**
make sure that all of the required directories exist
*/
func checkDirectories(p string) []string {
	configPath := path.Join(p, "config")
	dataPath := path.Join(p, "data")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println(fmt.Sprintf("Creating config directory: %s", configPath))
		os.Mkdir(configPath, os.ModeDir)
	}

	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		fmt.Println(fmt.Sprintf("Creating data directory: %s", dataPath))
		os.Mkdir(dataPath, os.ModeDir)
	}
	return []string{
		configPath,
		dataPath,
	}
}

/**
Create a default config and then read it (or existing file if it exists)
*/
func checkServerConfig(configPath string) (*serverConfig, error) {
	configFilePath := path.Join(configPath, "server.json")

	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		fmt.Println(fmt.Sprintf("Creating new config file: %s", configFilePath))

		// Create default config file...
		conf := serverConfig{
			ServerName: "Default MUD server",
			ServerPort: 1337,
			RoomData: struct {
				Config struct {
					MinWidth  int `json:"min_width"`
					MaxWidth  int `json:"max_width"`
					MinHeight int `json:"min_height"`
					MaxHeight int `json:"maxHeight"`
				} `json:"config"`
				MinRooms int `json:"min_rooms"`
			}{
				Config: struct {
					MinWidth  int `json:"min_width"`
					MaxWidth  int `json:"max_width"`
					MinHeight int `json:"min_height"`
					MaxHeight int `json:"maxHeight"`
				}{
					MinWidth:  5,
					MaxWidth:  100,
					MinHeight: 5,
					MaxHeight: 100,
				},
				MinRooms: 1,
			},
		}

		jData, _ := json.MarshalIndent(conf, "", " ")
		err := ioutil.WriteFile(configFilePath, jData, 0644)
		if err != nil {
			return nil, err
		}
		fmt.Println(fmt.Sprintf("Default config saved as %s", configFilePath))
	}

	if fData, err := ioutil.ReadFile(configFilePath); err != nil {
		return nil, err
	} else {
		fConfig := &serverConfig{}
		if err := json.Unmarshal(fData, fConfig); err != nil {
			return nil, err
		}
		return fConfig, nil
	}
}

func loadServerRooms(datapath string) error {
	roomFilePath := path.Join(datapath, "rooms.blob")
	_, err := os.Stat(roomFilePath)

	if os.IsNotExist(err) {

		for i := 0; i < 1; i++ {
			rand.Seed(time.Now().UnixNano())
			// Generate 1 room to start with
			width := rand.Intn(ServerInstance.config.RoomData.Config.MaxWidth-ServerInstance.config.RoomData.Config.MinWidth) + ServerInstance.config.RoomData.Config.MinWidth
			height := rand.Intn(ServerInstance.config.RoomData.Config.MaxHeight-ServerInstance.config.RoomData.Config.MinHeight) + ServerInstance.config.RoomData.Config.MinHeight
			fmt.Println(fmt.Sprintf("Generating a new room, size (width x height): %d x %d", width, height))

			newRoom := NewRoom(width, height)
			for x := 0; x < width; x++ {
				for y := 0; y < height; y++ {

					tt := TILE_TYPE_DIRT
					pass := true

					if x == 0 || y == 0 {
						tt = TILE_TYPE_WALL
						pass = false
					}

					t := Tile{
						Type:       tt,
						IsPassable: pass,
						Position:   Vector2{x, y},
					}

					newRoom.Tiles[x][y] = t
					newRoom.TileList = append(newRoom.TileList, t)
				}
			}
			ServerInstance.roomList[newRoom.ID.String()] = newRoom
			time.Sleep(time.Millisecond * 250)
		}
	}
	return nil
}
