package minecraft

import (
	"LaunchCore/eu.suro/launch/protos/server"
	"context"
	"time"

	"github.com/millkhan/mcstatusgo"
)

type routerServer struct {
	service Service
	server.UnimplementedServerServer
}

func NewRouterServer(service Service) server.ServerServer {
	return &routerServer{
		service: service,
	}
}

func (s *routerServer) CreateServer(ctx context.Context, req *server.CreateServerRequest) (res *server.Response, er error) {
	status, err := s.service.CreateServer(req.Version, req.Name, req.SaveWorld, req.Open)
	if err != nil {
		return nil, err
	}
	return &server.Response{
		Status: status,
	}, nil
}

func (s *routerServer) UpdateServer(ctx context.Context, req *server.UpdateServerRequest) (res *server.Response, er error) {
	status, err := s.service.UpdateServer(req.Port, req.Message)
	if err != nil {
		return nil, err
	}
	return &server.Response{
		Status: status,
	}, nil
}

func (s *routerServer) DeleteServer(ctx context.Context, req *server.DeleteServerRequest) (res *server.Response, er error) {
	status, err := s.service.DeleteServer(req.Port)
	if err != nil {
		return nil, err
	}
	return &server.Response{
		Status: status,
	}, nil
}

func (s *routerServer) ListServers(ctx context.Context, req *server.Empty) (res *server.ListServersResponse, er error) {
	servers, err := s.service.ListServers()
	if err != nil {
		return nil, err
	}
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
					Open:       servermc.Open,
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
			Open:       servermc.Open,
		})
	}
	return &server.ListServersResponse{
		Servers: srvInfo,
	}, nil
}
