package plugins

type Plugin struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	Name     string `json:"name"`
	SpigotID uint   `json:"spigot_id"`
}
