package user

import (
	"LaunchCore/internal/handlers"

	"github.com/gin-gonic/gin"
)

type handler struct {
}

func NewUserHandler() handlers.HandlerGroup {
	return &handler{}
}

func (h *handler) Register(r *gin.RouterGroup) {
	r.GET("/user", h.user)
}

func (h *handler) user(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.JSON(401, gin.H{
			"code":    401,
			"message": "Invalid token",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": user,
	})
}
