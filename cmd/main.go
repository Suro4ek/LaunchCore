package main

import (
	"LaunchCore/api/auth"
	"LaunchCore/api/middleware"
	"LaunchCore/internal/config"
	"LaunchCore/internal/minecraft"
	"LaunchCore/internal/minecraft/mc"
	"LaunchCore/internal/user"
	"LaunchCore/pkg/logging"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func main() {
	logging.Init()
	logging := logging.GetLogger()
	logging.Info("start application")
	cfg := config.GetConfig()
	//_ = mysql.NewClient(context.Background(), 3, cfg.MySQL)
	logging.Info("connect to MySQL")
	mc := startDocker(&logging)
	startWebServer(&logging, mc, &cfg.OAuth2)
}

func startWebServer(log *logging.Logger, MC *minecraft.MC, cfg *config.OAuth2) {
	router := gin.Default()
	log.Info("Start web server")
	conf := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes: []string{
			"user",
		},
		Endpoint: github.Endpoint,
	}

	auth := auth.NewAuthHandler(conf)
	auth.Register(router)

	api := router.Group("/api/")
	oauthm := middleware.NewOAuthMiddleware()
	api.Use(oauthm.CheckAuth())

	user := user.NewUserHandler()
	user.Register(api)

	mc := minecraft.NewMCHandler(minecraft.Deps{
		Log:       log,
		Client:    nil,
		Container: MC,
	})
	mc.Register(api)

	router.Run(":8080")
}

func startDocker(log *logging.Logger) *minecraft.MC {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(context.TODO(), "itzg/minecraft-server:latest", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	log.Info("pull image")
	io.Copy(log.Writer(), reader)
	defer reader.Close()
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}

	docker := mc.NewDocker(cli)
	return &docker
	//create timeout after 10 sec sleep
	//time.Sleep(10 * time.Second)
	//er = docker.Delete(id)
	//if er != nil {
	//	panic(er)
	//}
}
