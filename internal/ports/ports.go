package ports

import (
	"LaunchCore/pkg/logging"
	"LaunchCore/pkg/mysql"
)

type Ports struct {
	client *mysql.Client
	log    *logging.Logger
}

func NewPorts(client *mysql.Client, log *logging.Logger) *Ports {
	return &Ports{
		client: client,
		log:    log,
	}
}

//get available port
func (p *Ports) GetPort() uint32 {
	var port Port
	err := p.client.DB.Model(&Port{}).Where("used = ?", false).First(&port).Error
	if err != nil {
		return 0
	}
	port.Used = true
	p.log.Infof("get port %d", port.Port)
	p.client.DB.Save(&port)
	return port.Port
}

//free port
func (p *Ports) FreePort(port string) {
	var port1 Port
	err := p.client.DB.Model(&Port{}).Where("port = ?", port).First(&port1).Error
	if err != nil {
		return
	}
	p.log.Infof("free port %d", port1.Port)
	p.client.DB.Model(&Port{}).Where("port = ?", port).Update("used", false)
}
