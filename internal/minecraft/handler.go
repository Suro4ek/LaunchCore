package minecraft

import (
	"LaunchCore/internal/handlers"
	"LaunchCore/internal/version"
	"LaunchCore/pkg/logging"
	"LaunchCore/pkg/mysql"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/millkhan/mcstatusgo"
)

type handler struct {
	log    *logging.Logger
	client *mysql.Client
	mc     MC
}

type Deps struct {
	Log       *logging.Logger
	Client    *mysql.Client
	Container *MC
}

func NewMCHandler(deps Deps) handlers.HandlerGroup {
	return &handler{
		log:    deps.Log,
		client: deps.Client,
		mc:     *deps.Container,
	}
}

func (h *handler) Register(router *gin.RouterGroup) {
	minecraft := router.Group("minecraft")
	{
		minecraft.GET("/minecraft/servers", h.GetServers)
		minecraft.PATCH("/minecraft/servers/:id", h.UpdateServer)
		minecraft.POST("/minecraft/servers", h.CreateServer)
		minecraft.DELETE("/minecraft/servers/:id", h.DeleteServer)
	}
}

func (h *handler) GetServers(ctx *gin.Context) {
	var servers []Server
	h.client.DB.Model(&Server{}).Find(&servers)
	var srvInfo = make([]*ServerInfo, 0)
	for _, server := range servers {
		status, err := mcstatusgo.Status("0.0.0.0", server.Port, 10*time.Second, 5*time.Second)
		if err != nil {
			srvInfo = append(srvInfo, &ServerInfo{
				Players:    status.Players.Online,
				MaxPlayers: status.Players.Max,
				Version:    status.Version.Name,
				OwnerName:  server.OwnerName,
				Status:     "offline",
			})
			continue
		}
		srvInfo = append(srvInfo, &ServerInfo{
			Players:    status.Players.Online,
			MaxPlayers: status.Players.Max,
			Version:    status.Version.Name,
			OwnerName:  server.OwnerName,
			Status:     server.Status,
		})
	}
	ctx.JSON(200, srvInfo)
}

func (h *handler) UpdateServer(ctx *gin.Context) {
	err := h.client.DB.Model(&Server{}).Where("port = ?", ctx.Param("port")).Update("status", ctx.PostForm("status")).Error
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, "Server status succefull updated")
}

func (h *handler) CreateServer(ctx *gin.Context) {
	name := ctx.PostForm("name")
	port := ctx.PostForm("port")
	version_id := ctx.PostForm("version_id")
	var version version.Version
	err := h.client.DB.Where("id = ?", version_id).First(&version).Error
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	id, err := h.mc.Create(name, port, version.Name, version.JVVersion.String())
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	//string convert to uint16
	portUint16, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	h.client.DB.Create(&Server{
		Port:        uint16(portUint16),
		OwnerName:   name,
		ContainerID: id,
		Status:      "starting",
	})
	ctx.JSON(200, "Server created")
}

func (h *handler) DeleteServer(ctx *gin.Context) {
	var server Server
	err := h.client.DB.Where("port = ?", ctx.Param("port")).First(&server).Error
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = h.mc.Delete(server.ContainerID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = h.client.DB.Delete(server).Error
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, "Server deleted")
}
