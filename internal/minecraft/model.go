package minecraft

type Server struct {
	ID          uint   `json:"id" gorm:"primarykey"`
	Port        uint   `json:"port"`
	ContainerID string `json:"container_id"`
}

type ServerInfo struct {
	Players    int    `json:"players"`
	MaxPlayers int    `json:"max_players"`
	Version    string `json:"version"`
	Motd       string `json:"motd"`
}
