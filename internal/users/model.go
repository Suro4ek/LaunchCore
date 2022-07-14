package users

import "LaunchCore/internal/plugins"

type User struct {
	ID      uint              `json:"id" gorm:"primarykey"`
	Name    string            `json:"name"`
	Plugins []*plugins.Plugin `gorm:"many2many:user_plugins"`
	Friends []*User           `gorm:"many2many:user_friends"`
}

// type Friend struct {
// 	ID       uint `json:"id" gorm:"primarykey"`
// 	UserID   uint `json:"user_id"`
// 	FriendID uint `json:"friend_id"`
// 	Friend   User `gorm:"ForeignKey:FriendID;" json:"friend"`
// }
