package mc

import (
	"LaunchCore/internal/minecraft"
	"LaunchCore/internal/users"
	"LaunchCore/pkg/logging"
	"LaunchCore/pkg/mysql"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type docker struct {
	client *client.Client
	log    *logging.Logger
	DB     *mysql.Client
}

func NewDocker(client *client.Client, log *logging.Logger, db *mysql.Client) minecraft.MC {
	return &docker{
		client: client,
		log:    log,
		DB:     db,
	}
}

func (d *docker) Create(name string, port string, version string, java_version string, save_world bool, open bool) (id string, err error) {
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
	var mounts = make([]mount.Mount, 0)
	mounts = append(mounts, mount.Mount{
		Type:   mount.TypeBind,
		Source: "/plugins/",
		Target: "/data/plugins",
	})
	if save_world {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: "/data/minecraft/" + name + "/world",
			Target: "/data/world",
		})
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: "/data/minecraft/" + name + "/world_nether",
			Target: "/data/world_nether",
		})
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: "/data/minecraft/" + name + "/world_the_end",
			Target: "/data/world_the_end",
		})
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: "/data/minecraft/" + name + "/ops.json",
			Target: "/data/ops.json",
		})
	}
	var user users.User
	err = d.DB.DB.Model(&users.User{}).Where("name = ?", name).Association("Friends").Find(&user)
	env := []string{
		"SERVER_JAVA_OPTS=-Xmx2024M -Xms2024M -XX:+UseG1GC -XX:+UseStringDeduplication",
		"SERVER_MAX_PLAYERS=20",
		"EULA=TRUE",
		"USE_AIKAR_FLAGS=true",
		"ENABLE_AUTOSTOP=TRUE",
		"SYNC_SKIP_NEWER_IN_DESTINATION=false",
		"AUTOSTOP_TIMEOUT_EST=10",
		"OPS=" + name,
		"VERSION=" + version,
		"SPAWN_PROTECTION=0",
	}
	if !open {
		env = append(env, "ENFORCE_WHITELIST=TRUE")
		//array user.Friends to string
		var friends []string
		for _, friend := range user.Friends {
			friends = append(friends, friend.Friend.Name)
		}
		env = append(env, "WHITELIST="+strings.Join(friends, ",")+name)
	}
	if len(user.Plugins) > 0 {
		//int32 to string
		var plugins_str []string
		for _, plugin := range user.Plugins {
			plugins_str = append(plugins_str, fmt.Sprintf("%d", plugin))
		}
		env = append(env, "SPIGET_RESOURCES="+strings.Join(plugins_str, ","))
	}
	resp, err := d.client.ContainerCreate(context.Background(), &container.Config{
		Env:   env,
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
		Mounts: mounts,
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
