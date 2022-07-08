package mc

import (
	"LaunchCore/internal/minecraft"
	"LaunchCore/pkg/logging"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
)

type docker struct {
	client *client.Client
	log    *logging.Logger
}

func NewDocker(client *client.Client, log *logging.Logger) minecraft.MC {
	return &docker{
		client: client,
		log:    log,
	}
}

func (d *docker) Create(name string, port string, version string, java_version string) (id string, err error) {
	fmt.Println(name, port, version, java_version)
	reader, err := d.client.ImagePull(context.TODO(), "itzg/minecraft-server:"+java_version, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}
	_, err = io.Copy(d.log.Writer(), reader)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	resp, err := d.client.ContainerCreate(context.Background(), &container.Config{
		Env: []string{
			"SERVER_JAVA_OPTS=-Xmx2024M -Xms2024M -XX:+UseG1GC -XX:+UseStringDeduplication",
			"SERVER_MAX_PLAYERS=20",
			"EULA=TRUE",
			"USE_AIKAR_FLAGS=true",
			"AUTOSTOP_TIMEOUT_EST=10",
			"VERSION=" + version,
		},
		Image: "itzg/minecraft-server:" + java_version,
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
	}, nil, nil, "servermc-"+name)
	if err != nil {
		return "", err
	}

	err = d.client.ContainerStart(context.TODO(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}
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
