package mc

import (
	"LaunchCore/internal/minecraft"
	"LaunchCore/internal/users"
	"LaunchCore/pkg/logging"
	"LaunchCore/pkg/mysql"
	"context"
	"fmt"
	"io"
	"strconv"
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

func (d *docker) Create(name string, port int32, version string, java_version string, save_world bool, open bool) (id string, err error) {
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
	d.log.Info("pull image success")
	var user users.User
	err = d.DB.DB.Model(&users.User{}).Where("name = ?", name).Find(&user).Error
	if err != nil {
		return "", err
	}
	var mounts = make([]mount.Mount, 0)
	mounts = append(mounts, mount.Mount{
		Type:   mount.TypeBind,
		Source: "/root/server/",
		Target: "/config/",
	})
	mounts = append(mounts, mount.Mount{
		Type:   mount.TypeBind,
		Source: "/root/plugins/",
		Target: "/plugins",
	})
	if save_world {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: "/data/minecraft/" + user.Name + "/world",
			Target: "/data/world",
		})
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: "/data/minecraft/" + user.Name + "/world_nether",
			Target: "/data/world_nether",
		})
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: "/data/minecraft/" + user.Name + "/world_the_end",
			Target: "/data/world_the_end",
		})
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: "/data/minecraft/" + user.Name + "/ops.json",
			Target: "/data/ops.json",
		})
		d.log.Info("save world")
	}
	env := []string{
		"SERVER_JAVA_OPTS=-Xmx2024M -Xms2024M -XX:+UseG1GC -XX:+UseStringDeduplication",
		"SERVER_MAX_PLAYERS=20",
		"EULA=TRUE",
		"USE_AIKAR_FLAGS=true",
		"COPY_CONFIG_DEST=/data",
		"TYPE=PAPER",
		"SYNC_SKIP_NEWER_IN_DESTINATION=false",
		"SERVER_NAME=" + name,
		"VERSION=" + version,
		"ONLINE_MODE=FALSE",
		"SPAWN_PROTECTION=0",
		//AutoStop
		"ENABLE_AUTOSTOP=TRUE",
		"AUTOSTOP_TIMEOUT_INIT=100",
		"AUTOSTOP_TIMEOUT_EST=100",
		"OPS=" + user.RealName,
	}
	//if !open {
	//	env = append(env, "ENFORCE_WHITELIST=TRUE")
	//array user.Friends to string
	//var friends []string
	//for _, friend := range user.Friends {
	//	friends = append(friends, friend.Name)
	//}
	//env = append(env, "WHITELIST="+strings.Join(friends, ",")+name)
	//}
	if len(user.Plugins) > 0 {
		//int32 to string
		var plugins_str []string
		for _, plugin := range user.Plugins {
			plugins_str = append(plugins_str, plugin.SpigotID)
		}
		env = append(env, "SPIGET_RESOURCES="+strings.Join(plugins_str, ","))
		d.log.Info("plugins")
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
					HostPort: strconv.Itoa(int(port)) + "/tcp",
				},
			},
		},
		Mounts: mounts,
	}, nil, nil, "servermc-"+name)
	if err != nil {
		return "", err
	}
	d.log.Info("Created container: " + resp.ID)
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
	d.log.Info("Deleting container: " + id)
	err := d.client.ContainerRemove(context.TODO(), id, types.ContainerRemoveOptions{
		Force: true,
	})
	return err
}
