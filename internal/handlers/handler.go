package handlers

import "github.com/gin-gonic/gin"

type Handler interface {
	Register(router *gin.Engine)
}

type HandlerGroup interface {
	Register(router *gin.RouterGroup)
}
