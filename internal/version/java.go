package version

type JavaVersion uint

const (
	Undefined JavaVersion = iota
	java_17
	java_11
	java_8
)

func (jv JavaVersion) String() string {
	return [...]string{"latest", "java17", "java11", "java8-multiarch"}[jv]
}
