package minecraft

import (
	"LaunchCore/internal/plugins"
	"LaunchCore/internal/ports"
	"LaunchCore/internal/version"
	"LaunchCore/pkg/mysql"
	"errors"
)

type Service struct {
	ports  ports.Ports
	client *mysql.Client
	mc     MC
}

func NewMCService(ports ports.Ports, client *mysql.Client, mc MC) *Service {
	return &Service{
		ports:  ports,
		client: client,
		mc:     mc,
	}
}

func (s *Service) CreateServer(ver int32, name string, saveworld bool, open bool) (status string, err error) {
	var version version.Version
	err = s.client.DB.Where("id = ?", ver).First(&version).Error
	if err != nil {
		return "", err
	}
	//check server is exist db
	var server1 Server
	err = s.client.DB.Where("name = ?", name).First(&server1).Error
	if err == nil {
		return "", errors.New("server is exists")
	}
	port := s.ports.GetPort()
	if port == 0 {
		return "", errors.New("no free ports")
	}
	//port to string
	id, err := s.mc.Create(name, int32(port), version.Name, version.JVVersion, saveworld, open)
	if err != nil {
		return "", err
	}
	s.client.DB.Create(&Server{
		Port:        String(int32(port)),
		OwnerName:   name,
		ContainerID: id,
		Status:      "starting",
		Open:        open,
	})
	return "ok", nil
}

func (s *Service) UpdateServer(port int32, message string) (status string, err error) {
	err = s.client.DB.Model(&Server{}).Where("port = ?", port).Update("status", message).Error
	if err != nil {
		return "", err
	}
	return "ok", nil
}

func (s *Service) DeleteServer(port int32) (status string, err error) {
	var server1 Server
	err = s.client.DB.Where("port = ?", port).First(&server1).Error
	if err != nil {
		return "", err
	}
	err = s.mc.Delete(server1.ContainerID)
	if err != nil {
		return "", err
	}
	//value, err := strconv.ParseInt(server1.Port, 10, 32)
	s.ports.FreePort(server1.Port)
	err = s.client.DB.Delete(server1).Error
	if err != nil {
		return "", err
	}
	return "ok", nil
}

func (s *Service) ListServers() (server []Server, err error) {
	err = s.client.DB.Find(&server).Error
	if err != nil {
		return nil, err
	}
	return server, nil
}

func (s *Service) ListPlugins() (plugin []plugins.Plugin, err error) {
	err = s.client.DB.Find(&plugin).Error
	if err != nil {
		return nil, err
	}
	return plugin, nil
}

func (s *Service) ListVersions() (version []version.Version, err error) {
	err = s.client.DB.Find(&version).Error
	if err != nil {
		return nil, err
	}
	return version, nil
}

func String(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}
