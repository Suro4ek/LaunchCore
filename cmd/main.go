package main

import (
	"LaunchCore/eu.suro/launch/protos/server"
	"LaunchCore/eu.suro/launch/protos/user"
	"LaunchCore/eu.suro/launch/protos/web"
	"LaunchCore/internal/config"
	"LaunchCore/internal/minecraft"
	"LaunchCore/internal/minecraft/mc"
	"LaunchCore/internal/plugins"
	"LaunchCore/internal/ports"
	"LaunchCore/internal/users"
	"LaunchCore/internal/version"
	"LaunchCore/internal/webr"
	"LaunchCore/pkg/logging"
	"LaunchCore/pkg/mysql"
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/millkhan/mcstatusgo"
	"google.golang.org/grpc"
)

func main() {
	logging.Init()
	log := logging.GetLogger()
	log.Info("start application")
	cfg := config.GetConfig()
	newClient := mysql.NewClient(context.Background(), 3, cfg.MySQL)
	log.Info("connect to MySQL")
	docker := startDocker(&log, newClient)
	migrate(newClient)
	newPorts := ports.NewPorts(newClient, &log)
	checkServers(&log, newClient, docker, newPorts)
	deleteAllServerBeforeStart(newClient, *docker)
	startGRPCServer(&log, docker, newClient, *newPorts, cfg)
}

func migrate(client *mysql.Client) {
	client.DB.AutoMigrate(version.Version{})
	client.DB.AutoMigrate(minecraft.Server{})
	client.DB.AutoMigrate(plugins.Plugin{})
	client.DB.AutoMigrate(users.User{})
	client.DB.AutoMigrate(ports.Port{})
	// client.DB.AutoMigrate(users.Friend{})

}

func checkServers(log *logging.Logger, client *mysql.Client, mc *minecraft.MC, ports *ports.Ports) {
	var repeat = make([]uint32, 0)
	//var repeatBool = bool(false)
	ticker := time.NewTicker(2 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				var servers []minecraft.Server
				client.DB.Model(&minecraft.Server{}).Find(&servers)
				for _, servermc := range servers {
					value, err := strconv.ParseUint(servermc.Port, 10, 32)
					_, err = mcstatusgo.Status("0.0.0.0", uint16(value), 10*time.Second, 5*time.Second)
					if err != nil {
						if check(servermc.ID, repeat) {
							log.Infof("delete server %s by port uint16 %d int32 %d", servermc.OwnerName, uint16(value), int32(value))
							ports.FreePort(servermc.Port)
							(*mc).Delete(servermc.ContainerID)
							client.DB.Delete(servermc)
							delete(servermc.ID, repeat)
						} else {
							repeat = append(repeat, servermc.ID)
						}
					}
					continue
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func startGRPCServer(log *logging.Logger, mc *minecraft.MC, client *mysql.Client, ports ports.Ports, cfg *config.Config) {
	addr := "0.0.0.0:" + cfg.GRPCPort
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	service := minecraft.NewMCService(ports, client, *mc)
	router := minecraft.NewRouterServer(*service)
	webRouter := webr.NewWebRouter(*service, client)
	userRouter := users.NewRouterUser(client)
	log.Info("start grpc server")
	server.RegisterServerServer(s, router)
	user.RegisterUserServer(s, userRouter)
	web.RegisterWebServer(s, webRouter)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func deleteAllServerBeforeStart(client *mysql.Client, mc minecraft.MC) {
	var servers []minecraft.Server
	client.DB.Model(&minecraft.Server{}).Find(&servers)
	for _, servermc := range servers {
		client.DB.Delete(servermc)
		mc.Delete(servermc.ContainerID)
	}
}

func startDocker(log *logging.Logger, mysql *mysql.Client) *minecraft.MC {
	log.Info("start docker api")
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}

	docker := mc.NewDocker(cli, log, mysql)

	return &docker
}

//check value in []uint
func check(val uint32, list []uint32) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}

//delete value in []uint32
func delete(val uint32, list []uint32) []uint32 {
	for i, v := range list {
		if v == val {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}
