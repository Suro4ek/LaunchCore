package ports

type Port struct {
	ID   uint32 `json:"id" gorm:"primarykey"`
	Port uint32 `json:"port"`
	Used bool   `json:"used"`
}
