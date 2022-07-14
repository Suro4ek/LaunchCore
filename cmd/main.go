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
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/millkhan/mcstatusgo"
	"google.golang.org/grpc"
)

func main() {
	logging.Init()
	logging := logging.GetLogger()
	logging.Info("start application")
	cfg := config.GetConfig()
	client := mysql.NewClient(context.Background(), 3, cfg.MySQL)
	logging.Info("connect to MySQL")
	mc := startDocker(&logging, client)
	migrate(client)
	ports := ports.NewPorts(cfg.Minecraft.StartPort, cfg.Minecraft.EndPort)
	checkServers(client, mc, *ports)
	startGRPCServer(&logging, mc, client, *ports, cfg)
	deleteAllServerBeforeStart(client, *mc)
}

func migrate(client *mysql.Client) {
	client.DB.AutoMigrate(version.Version{})
	client.DB.AutoMigrate(minecraft.Server{})
	client.DB.AutoMigrate(plugins.Plugin{})
	client.DB.AutoMigrate(users.User{})
	// client.DB.AutoMigrate(users.Friend{})

}

func checkServers(client *mysql.Client, mc *minecraft.MC, ports ports.Ports) {
	var repeat = make([]uint, 0)
	var repeatBool = bool(false)
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
						if check(servermc.ID, repeat) {
							(*mc).Delete(servermc.ContainerID)
							ports.FreePort(int32(servermc.Port))
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
