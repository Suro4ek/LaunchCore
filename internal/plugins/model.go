package plugins

type Plugin struct {
	ID          uint   `json:"id" gorm:"primarykey"`
	Name        string `json:"name"`
	SpigotID    string `json:"spigot_id"`
	Description string `json:"description"`
}
