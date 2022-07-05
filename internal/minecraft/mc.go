package minecraft

type MC interface {
	Create(name string, port string) (id string, err error)
	Get(id string) error
	Delete(id string) error
}
