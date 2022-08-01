package minecraft

import "LaunchCore/internal/version"

type MC interface {
	Create(name string, port int32, version version.Version, save_world bool, open bool) (id string, err error)
	Get(id string) error
	Delete(id string) error
}
