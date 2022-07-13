package webr

import (
	"LaunchCore/eu.suro/launch/protos/web"
	"LaunchCore/internal/minecraft"
	"LaunchCore/internal/plugins"
	"LaunchCore/internal/version"
	"LaunchCore/pkg/mysql"
	"context"
)

type routerWeb struct {
	web.UnimplementedWebServer
	client  *mysql.Client
	service minecraft.Service
}

func NewWebRouter(service minecraft.Service, client *mysql.Client) web.WebServer {
	return &routerWeb{
		service: service,
		client:  client,
	}
}

func (r *routerWeb) DeleteAllServers(context.Context, *web.Empty) (*web.Response, error) {
	servers, err := r.service.ListServers()
	if err != nil {
		return nil, err
	}
	for _, server := range servers {
		r.service.DeleteServer(int32(server.Port))
	}
	return &web.Response{
		Message: "ok",
	}, nil
}

func (r *routerWeb) CreatePlugin(ctx context.Context, req *web.Plugin) (res *web.Response, err error) {
	err = r.client.DB.Model(plugins.Plugin{}).Create(&plugins.Plugin{
		Name:        req.Name,
		SpigotID:    uint(req.Spigotid),
		Description: req.Description,
	}).Error
	if err != nil {
		return nil, err
	}
	return &web.Response{
		Message: "ok",
	}, nil
}

func (r *routerWeb) DeletePlugin(ctx context.Context, req *web.Plugin) (res *web.Response, err error) {
	err = r.client.DB.Model(plugins.Plugin{}).Where("id = ?", req.Id).Delete(&plugins.Plugin{}).Error
	if err != nil {
		return nil, err
	}
	return &web.Response{
		Message: "ok",
	}, nil
}

func (r *routerWeb) GetPlugin(ctx context.Context, req *web.ResponseById) (res *web.Plugin, err error) {
	err = r.client.DB.Model(plugins.Plugin{}).Where("id = ?", req.Id).First(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (r *routerWeb) GetPlugins(ctx context.Context, req *web.Empty) (res *web.Plugins, err error) {
	err = r.client.DB.Model(plugins.Plugin{}).Find(&res.Plugins).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (r *routerWeb) UpdatePlugin(ctx context.Context, req *web.Plugin) (res *web.Response, err error) {
	err = r.client.DB.Model(plugins.Plugin{}).Where("id = ?", req.Id).Save(&plugins.Plugin{
		Name:        req.Name,
		SpigotID:    uint(req.Spigotid),
		Description: req.Description,
	}).Error
	if err != nil {
		return nil, err
	}
	return &web.Response{
		Message: "ok",
	}, nil
}
func (r *routerWeb) CreateVersion(ctx context.Context, req *web.Version) (res *web.Response, err error) {
	err = r.client.DB.Model(version.Version{}).Create(&version.Version{
		Name:        req.Name,
		Description: req.Description,
		Url:         req.Url,
		Version:     req.Version,
		JVVersion:   req.JavaVersion,
	}).Error
	if err != nil {
		return nil, err
	}
	return &web.Response{
		Message: "ok",
	}, nil
}
func (r *routerWeb) DeleteVersion(ctx context.Context, req *web.Version) (res *web.Response, err error) {
	err = r.client.DB.Model(version.Version{}).Where("id = ?", req.Id).Delete(&version.Version{}).Error
	if err != nil {
		return nil, err
	}
	return &web.Response{
		Message: "ok",
	}, nil
}
func (r *routerWeb) GetVersion(ctx context.Context, req *web.ResponseById) (res *web.Version, err error) {
	err = r.client.DB.Model(version.Version{}).Where("id = ?", req.Id).First(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (r *routerWeb) GetVersions(ctx context.Context, req *web.Empty) (res *web.Versions, err error) {
	err = r.client.DB.Model(version.Version{}).Find(&res.Versions).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (r *routerWeb) UpdateVersion(ctx context.Context, req *web.Version) (res *web.Response, err error) {
	err = r.client.DB.Model(version.Version{}).Where("id = ?", req.Id).Save(&version.Version{
		Name:        req.Name,
		Description: req.Description,
		Url:         req.Url,
		Version:     req.Version,
		JVVersion:   req.JavaVersion,
	}).Error
	if err != nil {
		return nil, err
	}
	return &web.Response{
		Message: "ok",
	}, nil
}
