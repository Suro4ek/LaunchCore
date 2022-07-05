package mc

import (
	"LaunchCore/internal/minecraft"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
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

func (d *docker) Create(name string, port string) (id string, err error) {
	resp, err := d.client.ContainerCreate(context.Background(), &container.Config{
		Env: []string{
			"SERVER_JAVA_OPTS=-Xmx1024M -Xms1024M -XX:+UseG1GC -XX:+UseStringDeduplication",
			"SERVER_MAX_PLAYERS=20",
			"EULA=TRUE",
			"USE_AIKAR_FLAGS=true",
			"AUTOSTOP_TIMEOUT_EST=3600",
			"VERSION=1.17",
		},
		Image: "itzg/minecraft-server",
	}, &container.HostConfig{
		AutoRemove: true,
		PortBindings: nat.PortMap{
			nat.Port("25565/tcp"): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port + "/tcp",
				},
			},
		},
	}, nil, nil, "container"+name)
	if err != nil {
		return "", err
	}
	fmt.Println(resp.ID)
	err = d.client.ContainerStart(context.TODO(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}
	info, er := d.client.ContainerInspect(context.TODO(), resp.ID)
	if er != nil {
		return "", er
	}
	fmt.Println(info.State.Status)
	return resp.ID, nil
}

func (d *docker) Get(id string) error {
	return nil
}

func (d *docker) Delete(id string) error {
	err := d.client.ContainerRemove(context.TODO(), id, types.ContainerRemoveOptions{
		Force: true,
	})
	return err
}
