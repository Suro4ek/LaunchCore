package webr

import (
	"LaunchCore/eu.suro/launch/protos/web"
	"LaunchCore/internal/minecraft"
	"LaunchCore/internal/plugins"
	"LaunchCore/internal/version"
	"LaunchCore/pkg/mysql"
	"context"
	"errors"
	"fmt"
	"strconv"
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

func (r *routerWeb) CreatePlugin(ctx context.Context, req *web.CreatePluginResponse) (res *web.Response, err error) {
	if req.GetName() == "" || req.GetDescription() == "" || req.GetSpigotid() == "" {
		return nil, errors.New("empty string")
	}
	err = r.client.DB.Model(plugins.Plugin{}).Create(&plugins.Plugin{
		Name:        req.Name,
		SpigotID:    req.Spigotid,
		Description: req.Description,
	}).Error
	if err != nil {
		return nil, err
	}
	return &web.Response{
		Message: "ok",
	}, nil
}

func (r *routerWeb) DeletePlugin(ctx context.Context, req *web.DeletePluginResponse) (res *web.Response, err error) {
	err = r.client.DB.Model(plugins.Plugin{}).Where("id = ?", req.Id).Delete(&plugins.Plugin{}).Error
	if err != nil {
		return nil, err
	}
	return &web.Response{
		Message: "ok",
	}, nil
}

func (r *routerWeb) GetPlugin(ctx context.Context, req *web.ResponseById) (res *web.Plugin, err error) {
	var plugin plugins.Plugin
	err = r.client.DB.Model(plugins.Plugin{}).Where("id = ?", req.Id).First(&plugin).Error
	if err != nil {
		return nil, err
	}
	var id string = strconv.FormatUint(uint64(plugin.ID), 10)
	res = &web.Plugin{
		Id:          id,
		Name:        plugin.Name,
		Spigotid:    plugin.SpigotID,
		Description: plugin.Description,
	}
	return res, nil
}
func (r *routerWeb) GetPlugins(ctx context.Context, req *web.Empty) (res *web.Plugins, err error) {
	var plugins1 []plugins.Plugin
	err = r.client.DB.Model(plugins.Plugin{}).Find(&plugins1).Error
	if err != nil {
		return nil, err
	}
	var pluginsWeb []*web.Plugin

	for _, plugin := range plugins1 {
		var id string = strconv.FormatUint(uint64(plugin.ID), 10)
		pluginsWeb = append(pluginsWeb, &web.Plugin{
			Id:          id,
			Name:        plugin.Name,
			Description: plugin.Description,
			Spigotid:    plugin.SpigotID,
		})
	}
	res = &web.Plugins{
		Plugins: pluginsWeb,
	}
	return res, nil
}
func (r *routerWeb) UpdatePlugin(ctx context.Context, req *web.Plugin) (res *web.Response, err error) {
	fmt.Println(req.Id)
	var plugin plugins.Plugin
	err = r.client.DB.Model(plugins.Plugin{}).Where("id = ?", req.Id).Find(&plugin).Error
	if err != nil {
		return nil, err
	}
	if req.GetName() == "" {
		req.Name = plugin.Name
	}
	if req.GetDescription() == "" {
		req.Description = plugin.Description
	}
	if req.GetSpigotid() == "" {
		req.Spigotid = plugin.SpigotID
	}
	id, err := strconv.ParseUint(req.Id, 10, 64)
	err = r.client.DB.Save(&plugins.Plugin{
		ID:          uint(id),
		Name:        req.Name,
		SpigotID:    req.Spigotid,
		Description: req.Description,
	}).Error
	if err != nil {
		return nil, err
	}
	return &web.Response{
		Message: "ok",
	}, nil
}
func (r *routerWeb) CreateVersion(ctx context.Context, req *web.CreateVersionResponse) (res *web.Response, err error) {
	if req.GetName() == "" || req.GetVersion() == "" || req.GetJavaVersion() == "" || req.GetDescription() == "" {
		return nil, errors.New("empty string")
	}
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
func (r *routerWeb) DeleteVersion(ctx context.Context, req *web.DeleteVersionResponse) (res *web.Response, err error) {
	err = r.client.DB.Model(version.Version{}).Where("id = ?", req.Id).Delete(&version.Version{}).Error
	if err != nil {
		return nil, err
	}
	return &web.Response{
		Message: "ok",
	}, nil
}
func (r *routerWeb) GetVersion(ctx context.Context, req *web.ResponseById) (res *web.Version, err error) {
	var version1 version.Version
	err = r.client.DB.Model(version.Version{}).Where("id = ?", req.Id).First(&version1).Error
	if err != nil {
		return nil, err
	}
	var id string = strconv.FormatUint(uint64(version1.ID), 10)
	res = &web.Version{
		Id:          id,
		Name:        version1.Name,
		Description: version1.Description,
		Url:         version1.Url,
		Version:     version1.Version,
		JavaVersion: version1.JVVersion,
	}
	return res, nil
}

func (r *routerWeb) GetVersions(ctx context.Context, req *web.Empty) (res *web.Versions, err error) {
	var versions []version.Version
	err = r.client.DB.Model(version.Version{}).Find(&versions).Error
	if err != nil {
		return nil, err
	}
	var versionsWeb []*web.Version
	for _, version := range versions {
		var id string = strconv.FormatUint(uint64(version.ID), 10)
		versionsWeb = append(versionsWeb, &web.Version{
			Id:          id,
			Name:        version.Name,
			Description: version.Description,
			Url:         version.Url,
			Version:     version.Version,
			JavaVersion: version.JVVersion,
		})
	}
	res = &web.Versions{
		Versions: versionsWeb,
	}
	return res, nil
}

func (r *routerWeb) UpdateVersion(ctx context.Context, req *web.Version) (res *web.Response, err error) {
	var ver version.Version
	err = r.client.DB.Model(version.Version{}).Where("id = ?", req.Id).Find(&ver).Error
	if err != nil {
		return nil, err
	}
	if req.GetVersion() == "" {
		req.Version = ver.Version
	}
	if req.GetName() == "" {
		req.Name = ver.Name
	}
	if req.GetDescription() == "" {
		req.Description = ver.Description
	}
	if req.GetJavaVersion() == "" {
		req.JavaVersion = ver.JVVersion
	}
	id, err := strconv.ParseUint(req.Id, 10, 64)
	err = r.client.DB.Save(&version.Version{
		ID:          uint(id),
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
