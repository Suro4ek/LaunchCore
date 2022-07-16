package minecraft

type MC interface {
	Create(name string, port int32, version string, java_version string, save_world bool, open bool) (id string, err error)
	Get(id string) error
	Delete(id string) error
}
