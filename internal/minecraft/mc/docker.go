package mc

import (
	"LaunchCore/internal/minecraft"
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type docker struct {
	client *client.Client
}

func NewDocker(client *client.Client) minecraft.MC {
	return &docker{
		client: client,
	}
}

func (d *docker) Create(name string, port string) error {
	d.client.ContainerCreate(context.Background(), &container.Config{
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
	return nil
}
func (d *docker) Get(id string) error {
	return nil
}
func (d *docker) Delete(id string) error {
	return nil
}
