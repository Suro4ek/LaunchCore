package minecraft

type Server struct {
	ID          uint32 `json:"id" gorm:"primarykey"`
	Port        string `json:"port"`
	ContainerID string `json:"container_id"`
	OwnerName   string `json:"owner_name"`
	Status      string `json:"status"`
	Open        bool   `json:"open"`
	Type        string `json:"type"`
}

type ServerInfo struct {
	Players    int    `json:"players"`
	MaxPlayers int    `json:"max_players"`
	Version    string `json:"version"`
	OwnerName  string `json:"owner_name"`
	Status     string `json:"status"`
	Port       string `json:"port"`
}
