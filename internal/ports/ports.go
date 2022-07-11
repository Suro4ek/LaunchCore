package ports

type Ports struct {
	startPort int32
	endPort   int32
	Ports     []int32
	UsedPorts []int32
}

func NewPorts(startPort int32, EndPort int32) *Ports {
	ports := &Ports{
		startPort: startPort,
		endPort:   EndPort,
		Ports:     make([]int32, 0),
		UsedPorts: make([]int32, 0),
	}
	for i := startPort; i <= EndPort; i++ {
		ports.Ports = append(ports.Ports, i)
	}
	return ports
}

//get available port
func (p *Ports) GetPort() int32 {
	for _, port := range p.Ports {
		if !p.IsUsed(port) {
			p.UsedPorts = append(p.UsedPorts, port)
			return port
		}
	}
	return 0
}

//isUsed port
func (p *Ports) IsUsed(port int32) bool {
	for _, usedPort := range p.UsedPorts {
		if usedPort == port {
			return true
		}
	}
	return false
}

//free port
func (p *Ports) FreePort(port int32) {
	for i, usedPort := range p.UsedPorts {
		if usedPort == port {
			p.UsedPorts = append(p.UsedPorts[:i], p.UsedPorts[i+1:]...)
			return
		}
	}
}
