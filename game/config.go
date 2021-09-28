package game

type serverConfig struct {
	ServerName      string `json:"server_name"`
	ServerPort      int    `json:"server_port"`
	ServerAddress   string `json:"server_address"`
	PingConnections bool   `json:"ping_connections"`
	RoomData        struct {
		Config struct {
			MinWidth  int `json:"min_width"`
			MaxWidth  int `json:"max_width"`
			MinHeight int `json:"min_height"`
			MaxHeight int `json:"maxHeight"`
		} `json:"config"`
		MinRooms int `json:"min_rooms"`
	} `json:"room_data"`
}
