package minecraft

import (
	"LaunchCore/eu.suro/launch/protos/server"
	"LaunchCore/internal/version"
	"LaunchCore/pkg/mysql"
	"context"
	"time"

	"github.com/millkhan/mcstatusgo"
)

type routerServer struct {
	ports  []int
	client *mysql.Client
	server.UnimplementedServerServer
	mc MC
}

type Deps struct {
	Client *mysql.Client
	Mc     MC
	Ports  []int
}

func NewRouterServer(deps Deps) server.ServerServer {
	return &routerServer{
		client: deps.Client,
		mc:     deps.Mc,
		ports:  deps.Ports,
	}
}

func (s *routerServer) CreateServer(ctx context.Context, req *server.CreateServerRequest) (res *server.Response, er error) {
	var version version.Version
	err := s.client.DB.Where("id = ?", req.GetVersion()).First(&version).Error
	if err != nil {
		return nil, err
	}
	//int32 to string
	id, err := s.mc.Create(req.Name, String(req.Port), version.Name, version.JVVersion.String())
	if err != nil {

		return nil, err
	}
	//string convert to uint16
	s.client.DB.Create(&Server{
		Port:        uint16(req.Port),
		OwnerName:   req.Name,
		ContainerID: id,
		Status:      "starting",
	})
	return &server.Response{
		Status: "ok",
	}, nil
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

func (s *routerServer) UpdateServer(ctx context.Context, req *server.UpdateServerRequest) (res *server.Response, er error) {
	err := s.client.DB.Model(&Server{}).Where("port = ?", req.Port).Update("status", req.Message).Error
	if err != nil {
		return nil, err
	}
	return &server.Response{
		Status: "ok",
	}, nil
}

func (s *routerServer) DeleteServer(ctx context.Context, req *server.DeleteServerRequest) (res *server.Response, er error) {
	var server1 Server
	err := s.client.DB.Where("port = ?", req.Port).First(&server1).Error
	if err != nil {
		return nil, err
	}
	err = s.mc.Delete(server1.ContainerID)
	if err != nil {
		return nil, err
	}
	err = s.client.DB.Delete(server1).Error
	if err != nil {
		return nil, err
	}
	return &server.Response{
		Status: "ok",
	}, nil
}

func (s *routerServer) ListServers(ctx context.Context, req *server.Empty) (res *server.ListServersResponse, er error) {
	var servers []Server
	s.client.DB.Model(&Server{}).Find(&servers)
	var srvInfo = make([]*server.ServerInfo, 0)
	for _, servermc := range servers {
		status, err := mcstatusgo.Status("0.0.0.0", servermc.Port, 10*time.Second, 5*time.Second)
		if err != nil {
			if servermc.Status == "starting" {
				srvInfo = append(srvInfo, &server.ServerInfo{
					Players:    int32(status.Players.Online),
					Maxplayers: int32(status.Players.Max),
					Version:    status.Version.Name,
					OwnerName:  servermc.OwnerName,
					Status:     "starting",
				})
			}
			continue
		}
		srvInfo = append(srvInfo, &server.ServerInfo{
			Players:    int32(status.Players.Online),
			Maxplayers: int32(status.Players.Max),
			Version:    status.Version.Name,
			OwnerName:  servermc.OwnerName,
			Status:     servermc.Status,
		})
	}
	return &server.ListServersResponse{
		Servers: srvInfo,
	}, nil
}
