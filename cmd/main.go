package main

import (
	"LaunchCore/pkg/logging"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func main() {
	logging.Init()
	logging := logging.GetLogger()
	logging.Info("start application")
	// cfg := config.GetConfig()
	// _ = mysql.NewClient(context.Background(), 3, cfg.MySQL)
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

	reader, err := cli.ImagePull(context.TODO(), "itzg/minecraft-server:java8", types.ImagePullOptions{})
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
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Env: []string{
			"SERVER_JAVA_OPTS=-Xmx1024M -Xms1024M -XX:+UseG1GC -XX:+UseStringDeduplication",
			"SERVER_MAX_PLAYERS=20",
			"EULA=TRUE",
			"SERVER_NAME=LaunchCore",
			"VERSION=1.15.2",
		},
		ExposedPorts: nat.PortSet{
			nat.Port("25565/tcp"): struct{}{},
		},
		Image: "itzg/minecraft-server:java8",
		Cmd:   []string{"echo", "hello world"},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"25565/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "25565",
				},
			},
		},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	cli.ContainerStart(context.TODO(), resp.ID, types.ContainerStartOptions{})
	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}
