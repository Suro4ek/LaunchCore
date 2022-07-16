package version

type Version struct {
	ID          uint32 `json:"id" gorm:"primarykey"`
	JVVersion   string `json:"java_version"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Version     string `json:"version"`
}
