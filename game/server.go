package game

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/Entrio/subenv"
	"github.com/rs/zerolog/log"
)

var ServerInstance *Server

type (
	Server struct {
		gameLoop        gameLoop
		connectionsList []*Connection
		roomList        map[string]*Room
		ticker          *time.Ticker
		config          *serverConfig
		packetHandler   map[PacketType]PacketHandler
	}
	gameLoop struct {
		ticker *time.Ticker
	}
)

// GetServer returns an existing instance or creates a new one /**
func GetServer() (*Server, error) {
	if ServerInstance == nil {
		log.Debug().Msg("No server instance initialized, creating a new one...")
		handlers := make(map[PacketType]PacketHandler)

		handlers[MsgUpdateRoomPayload] = RoomUpdateHandler{}
		handlers[MsgRoomCountRequest] = RoomCountHandler{}

		log.Debug().Int("count", len(handlers)).Msg("Total handlers")

		ServerInstance = &Server{
			connectionsList: make([]*Connection, 0),
			roomList:        map[string]*Room{},
			packetHandler:   handlers,
		}
	}

	// Get current working directory
	cwd, _ := os.Getwd()
	log.Debug().Str("cwd", cwd).Msg("Current directory")

	dirs := checkDirectories(cwd)
	config, err := checkServerConfig(dirs[0])

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to check directories")
		return nil, err
	}
	ServerInstance.config = config

	/*
		err = loadServerRooms(dirs[1])
		if err != nil {
			return nil, err
		}
	*/
	return ServerInstance, nil
}

// Start creates a new ticked and launched a goroutine the periodically pings clients and performs game loop
func (server *Server) Start() {
	log.Info().Msg("Starting server game loop and ping goroutines")
	server.ticker = time.NewTicker(time.Millisecond * 3000)
	if &server.gameLoop == nil {
		log.Debug().Msg("Creating new game loop")
		server.gameLoop = gameLoop{
			ticker: time.NewTicker(time.Millisecond * 33),
		}
	}
	go func() {
		log.Info().Msg("Starting game loop")
		for range server.gameLoop.ticker.C {
			// Do game loop login and handlers here
		}
	}()

	go func() {
		log.Debug().Msg("Starting ping goroutine")
		for range server.ticker.C {
			pkt := NewPacket(MsgPingRequest)

			for _, c := range server.connectionsList {
				// game loop
				if server.config.PingConnections {
					sendMessageToConnection(c, *pkt)
				}
			}

		}
	}()
}

func (server *Server) GetPort() int {
	return server.config.ServerPort
}

func (server *Server) GetName() string {
	return server.config.ServerName
}

// FindRoom fetches a room based on UUID string
func (server *Server) FindRoom(uuid string) *Room {

	if r, found := server.roomList[uuid]; found {
		return r
	}

	return nil
}

// checkDirectories makes sure that all of the required directories exist
func checkDirectories(cwd string) []string {
	var err error
	configPath := path.Join(cwd, "config")
	dataPath := path.Join(cwd, "data")
	log.Trace().Str("data", dataPath).Str("config", configPath).Msg("Setting directories")

	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		log.Debug().Str("config", configPath).Msg("Config path does not exist, creating...")
		err = os.Mkdir(configPath, 0777)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to create config path")
		}
	}

	if _, err = os.Stat(dataPath); os.IsNotExist(err) {
		log.Debug().Str("data", dataPath).Msg("Data path does not exist, creating...")
		err = os.Mkdir(dataPath, 0777)

		if err != nil {
			log.Warn().Err(err).Msg("Failed to create data path")
		}
	}
	log.Debug().Msg("Created data and config paths successfully")
	return []string{
		configPath,
		dataPath,
	}
}

// checkServerConfig creates a default config and then read it (or existing file if it exists)
func checkServerConfig(configPath string) (*serverConfig, error) {
	var err error
	configFilePath := path.Join(configPath, subenv.Env("SERVER_CONFIG", "server.json"))
	log.Trace().Str("server_config", subenv.Env("SERVER_CONFIG", "server.json")).Msg("Server config filename, use 'SERVER_CONFIG' environment variable to overwrite")

	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		log.Info().Msg("Server config file does not exist, creating a new one")

		// Create default config file...
		conf := serverConfig{
			ServerName:      "Default MUD server",
			ServerPort:      1337,
			PingConnections: false,
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
				MinRooms: 6,
			},
		}

		jData, _ := json.MarshalIndent(conf, "", " ")
		err = ioutil.WriteFile(configFilePath, jData, 0666)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to write config file to disk")
		}
		log.Info().Str("filepath", configFilePath).Msg("Created config file successfully")
		log.Trace().Str("data", string(jData)).Msg("Server config")
	}

	if fData, err := ioutil.ReadFile(configFilePath); err != nil {
		log.Warn().Err(err).Msg("Failed to read config file from disk")
		return nil, err
	} else {
		fConfig := &serverConfig{}
		if err := json.Unmarshal(fData, fConfig); err != nil {
			log.Warn().Err(err).Msg("failed to unmarshal json data")
			return nil, err
		}
		log.Debug().Msg("Returning server config")
		return fConfig, nil
	}
}

func loadServerRooms(dataPath string) error {
	roomFilePath := path.Join(dataPath, "rooms.blob")
	_, err := os.Stat(roomFilePath)

	if os.IsNotExist(err) {

		for i := 0; i < ServerInstance.config.RoomData.MinRooms; i++ {
			time.Sleep(time.Millisecond * 250)
			rand.Seed(time.Now().UnixNano())
			// Generate 1 room to start with
			width := rand.Intn(ServerInstance.config.RoomData.Config.MaxWidth-ServerInstance.config.RoomData.Config.MinWidth) + ServerInstance.config.RoomData.Config.MinWidth
			height := rand.Intn(ServerInstance.config.RoomData.Config.MaxHeight-ServerInstance.config.RoomData.Config.MinHeight) + ServerInstance.config.RoomData.Config.MinHeight
			fmt.Println(fmt.Sprintf("Generating a new room, size (width x height): %d x %d", width, height))

			newRoom := NewRoom(width, height)
			for x := 0; x < width; x++ {
				for y := 0; y < height; y++ {

					tt := TILE_TYPE_AIR
					pass := false

					if x == 0 || y == 0 {
						tt = TILE_TYPE_WALL
						pass = false
					}

					if x+1 == width || y+1 == height {
						tt = TILE_TYPE_WALL
						pass = false
					}

					t := Tile{
						Type:       tt,
						IsPassable: pass,
						Position:   Vector2{x, y},
					}

					newRoom.Tiles[x][y] = t
				}
			}
			ServerInstance.roomList[newRoom.ID.String()] = newRoom
		}
		go saveServerRooms(dataPath)
	}
	return nil
}

/**
Room structure:
uint16 number of rooms
.... cycle of rooms
36 bytes - room uuid
uint16 - number of bytes used for description (max 65535)
.. <n> bytes - description
byte - 1 byte used to indicate if room is active or not
*/
func saveServerRooms(dataPath string) error {
	roomFilePath := path.Join(dataPath, "rooms.blob")
	_, err := os.Stat(roomFilePath)

	if os.IsNotExist(err) {
		// no room index
		buffer := make([]byte, 0)

		tBuffer := make([]byte, 2)
		binary.LittleEndian.PutUint16(tBuffer, uint16(len(ServerInstance.roomList)))

		buffer = append(buffer, tBuffer...)

		for _, v := range ServerInstance.roomList {
			buffer = append(buffer, []byte(v.ID.String())...)

			dBuffer := make([]byte, 2)
			binary.LittleEndian.PutUint16(dBuffer, uint16(len(v.Description)))
			buffer = append(buffer, dBuffer...)
			buffer = append(buffer, []byte(v.Description)...)

			if v.isActive {
				buffer = append(buffer, 0)
			} else {
				buffer = append(buffer, 1)
			}
		}
		err := ioutil.WriteFile(roomFilePath, buffer, 0644)
		check(err)

	} else {
		// this file exists, need to update it
	}
	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
