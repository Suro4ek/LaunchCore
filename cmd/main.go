package main

import (
	"LaunchCore/internal/minecraft/mc"
	"LaunchCore/pkg/logging"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
)

func main() {
	logging.Init()
	logging := logging.GetLogger()
	logging.Info("start application")
	//cfg := config.GetConfig()
	//_ = mysql.NewClient(context.Background(), 3, cfg.MySQL)
	logging.Info("connect to MySQL")
	startWebServer()
	startDocker(&logging)
}

func startWebServer() {

}

func startDocker(log *logging.Logger) {
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
	_, er := docker.Create("test", "25567")
	if er != nil {
		panic(er)
	}
	//create timeout after 10 sec sleep
	//time.Sleep(10 * time.Second)
	//er = docker.Delete(id)
	//if er != nil {
	//	panic(er)
	//}
}
