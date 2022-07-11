package users

type User struct {
	ID      uint     `json:"id" gorm:"primarykey"`
	Name    string   `json:"name"`
	Plugins []int32  `json:"plugins"`
	Friends []Friend `gorm:"ForeignKey:UserID;" json:"friend"`
}

type Friend struct {
	ID       uint `json:"id" gorm:"primarykey"`
	UserID   uint `json:"user_id"`
	FriendID uint `json:"friend_id"`
	Friend   User `gorm:"ForeignKey:FriendID;" json:"friend"`
}
