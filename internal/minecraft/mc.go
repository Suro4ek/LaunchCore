package minecraft

type MC interface {
	Create(name string, port string, version string, java_version string, save_world bool) (id string, err error)
	Get(id string) error
	Delete(id string) error
}
