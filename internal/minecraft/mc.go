package minecraft

type MC interface {
	Create(name string, port string) error
	Get(id string) error
	Delete(id string) error
}
