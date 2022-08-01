package mc

import (
	"LaunchCore/internal/minecraft"
	"LaunchCore/internal/users"
	"LaunchCore/internal/version"
	"LaunchCore/pkg/logging"
	"LaunchCore/pkg/mysql"
	"context"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"gorm.io/gorm"
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

func (d *docker) Create(name string, port int32, version version.Version, save_world bool, open bool) (id string, err error) {
	reader, err := d.client.ImagePull(context.TODO(), "itzg/minecraft-server:"+version.JVVersion, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}
	_, err = io.Copy(d.log.Writer(), reader)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	d.log.Info("pull image success")
	err = d.DB.DB.Model(&minecraft.Server{}).Where("owner_name = ?", name).First(&minecraft.Server{}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}
	} else {
		return "", errors.New("server already exists")
	}
	var user users.User
	err = d.DB.DB.Model(&users.User{}).Where("name = ?", name).Find(&user).Error
	if err != nil {
		return "", err
	}
	path, _ := os.Getwd()
	var mounts = make([]mount.Mount, 0)
	os.Mkdir(path+"/config/", os.ModePerm)
	mounts = append(mounts, mount.Mount{
		Type:   mount.TypeBind,
		Source: path + "/config/",
		Target: "/config/",
	})
	os.Mkdir(path+"/plugins/", os.ModePerm)
	mounts = append(mounts, mount.Mount{
		Type:   mount.TypeBind,
		Source: path + "/plugins/",
		Target: "/plugins",
	})
	if save_world && version.Url == "" {
		os.MkdirAll(path+"/worlds/"+user.Name+"/world", os.ModePerm)
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: path + "/worlds/" + user.Name + "/world/",
			Target: "/data/world/",
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
		"UID=1001",
		"GID=1001",
		"SYNC_SKIP_NEWER_IN_DESTINATION=false",
		"SERVER_NAME=" + name,
		"VERSION=" + version.Version,
		"ONLINE_MODE=FALSE",
		"SPAWN_PROTECTION=0",
		//AutoStop
		"ENABLE_AUTOSTOP=TRUE",
		"AUTOSTOP_TIMEOUT_INIT=60",
		"AUTOSTOP_TIMEOUT_EST=60",
		"OPS=" + user.RealName,
	}
	if version.Url != "" {
		env = append(env, "WORLD="+version.Url)
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
		Image: "itzg/minecraft-server:" + version.JVVersion,
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
