package minecraft

import (
	"LaunchCore/eu.suro/launch/protos/server"
	"context"
	"fmt"
	"strconv"
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
		value, err := strconv.ParseUint(servermc.Port, 10, 16)
		status, err := mcstatusgo.Status("0.0.0.0", uint16(value), 10*time.Second, 5*time.Second)
		if err != nil {
			if servermc.Status == "starting" {
				srvInfo = append(srvInfo, &server.ServerInfo{
					Players:    int32(status.Players.Online),
					Maxplayers: int32(status.Players.Max),
					Version:    status.Version.Name,
					OwnerName:  servermc.OwnerName,
					Status:     "starting",
					Open:       servermc.Open,
					Port:       servermc.Port,
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
			Port:       servermc.Port,
		})
	}
	return &server.ListServersResponse{
		Servers: srvInfo,
	}, nil
}

func (r *routerServer) GetPlugins(context.Context, *server.Empty) (*server.Plugins, error) {
	plugins1, err := r.service.ListPlugins()
	var plugins = make([]*server.Plugin, 0)
	for _, plugin := range plugins1 {
		plugins = append(plugins, &server.Plugin{
			Id:          plugin.ID,
			Name:        plugin.Name,
			Spigotid:    plugin.SpigotID,
			Description: plugin.Description,
		},
		)
	}
	if err != nil {
		return nil, err
	}
	return &server.Plugins{
		Plugins: plugins,
	}, nil
}

func (r *routerServer) GetVersions(context.Context, *server.Empty) (*server.Versions, error) {
	versions1, err := r.service.ListVersions()
	var versions = make([]*server.Version, 0)
	for _, version := range versions1 {
		fmt.Println(version.ID, version.Name)
		versions = append(versions, &server.Version{
			Id:          version.ID,
			Name:        version.Name,
			Description: version.Description,
			Url:         version.Url,
			Version:     version.Version,
			JavaVersion: version.JVVersion,
		},
		)
	}
	if err != nil {
		return nil, err
	}
	fmt.Println(versions[0].Id, versions[0].Name)
	return &server.Versions{
		Versions: versions,
	}, nil
}

func (r *routerServer) DeleteAllServers(context.Context, *server.Empty) (*server.Empty, error) {
	servers, err := r.service.ListServers()
	if err != nil {
		return nil, err
	}
	for _, server := range servers {
		value, err := strconv.ParseInt(server.Port, 10, 32)
		if err != nil {
			return nil, err
		}
		r.service.DeleteServer(int32(value))
	}
	return &server.Empty{}, nil
}
