package minecraft

type Server struct {
	ID          uint   `json:"id" gorm:"primarykey"`
	Port        uint16 `json:"port"`
	ContainerID string `json:"container_id"`
	OwnerName   string `json:"owner_name"`
	Status      string `json:"status"`
}

type ServerInfo struct {
	Players    int    `json:"players"`
	MaxPlayers int    `json:"max_players"`
	Version    string `json:"version"`
	OwnerName  string `json:"owner_name"`
	Status     string `json:"status"`
}
