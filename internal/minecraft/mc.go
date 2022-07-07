package minecraft

type MC interface {
	Create(name string, port string, version string, java_version string) (id string, err error)
	Get(id string) error
	Delete(id string) error
}
