package minecraft

import (
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

func (s *Service) CreateServer(ver string, name string, saveworld bool, open bool) (status string, err error) {
	var version version.Version
	err = s.client.DB.Where("id = ?", ver).First(&version).Error
	if err != nil {
		return "", err
	}
	port := s.ports.GetPort()
	if port == 0 {
		return "", errors.New("no free ports")
	}
	//port to string
	id, err := s.mc.Create(name, string(port), version.Name, version.JVVersion, saveworld, open)
	if err != nil {
		return "", err
	}
	s.client.DB.Create(&Server{
		Port:        uint16(port),
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
	s.ports.FreePort(int32(server1.Port))
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
