package main

import (
	"LaunchCore/eu.suro/launch/protos/server"
	"LaunchCore/internal/config"
	"LaunchCore/internal/minecraft"
	"LaunchCore/internal/minecraft/mc"
	"LaunchCore/pkg/logging"
	"LaunchCore/pkg/mysql"
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/millkhan/mcstatusgo"
	"google.golang.org/grpc"
)

var repeat = make([]uint, 0)
var repeatBool = bool(false)

func main() {
	logging.Init()
	logging := logging.GetLogger()
	logging.Info("start application")
	cfg := config.GetConfig()
	client := mysql.NewClient(context.Background(), 3, cfg.MySQL)
	logging.Info("connect to MySQL")
	mc := startDocker(&logging)
	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				var servers []minecraft.Server
				client.DB.Model(&minecraft.Server{}).Find(&servers)
				for _, servermc := range servers {
					_, err := mcstatusgo.Status("0.0.0.0", servermc.Port, 10*time.Second, 5*time.Second)
					if err != nil {
						if check(servermc.ID, repeat) && servermc.Status != "starting" {
							(*mc).Delete(servermc.ContainerID)
							client.DB.Delete(servermc)
							delete(servermc.ID, repeat)
							continue
						}
						repeat = append(repeat, servermc.ID)
						continue
					}
					if check(servermc.ID, repeat) {
						delete(servermc.ID, repeat)
					}
				}
				if repeatBool && len(repeat) > 0 {
					repeat = make([]uint, 0)
					repeatBool = false
				} else {
					repeatBool = true
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	startGRPCServer(&logging, mc)
}

func startGRPCServer(log *logging.Logger, mc *minecraft.MC) {
	addr := "0.0.0.0:" + "9000"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//create new grpc server and register its dependencies
	s := grpc.NewServer()
	router := minecraft.NewRouterServer(minecraft.Deps{
		Client: nil,
		Mc:     *mc,
	})
	log.Info("start grpc server")
	server.RegisterServerServer(s, router)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func startDocker(log *logging.Logger) *minecraft.MC {
	log.Info("start docker api")
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
}

//check value in []uint
func check(val uint, list []uint) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}

//delete value in []uint
func delete(val uint, list []uint) []uint {
	for i, v := range list {
		if v == val {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}
