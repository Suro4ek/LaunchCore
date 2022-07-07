package auth

import (
	"LaunchCore/internal/handlers"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/xsrftoken"
	"golang.org/x/oauth2"
)

type handler struct {
	Conf *oauth2.Config
}

func NewAuthHandler(conf *oauth2.Config) handlers.Handler {
	return &handler{
		Conf: conf,
	}
}

func (h *handler) Register(router *gin.Engine) {
	router.GET("/auth/login", h.login)
	router.GET("/auth/oauth", h.oauth)
}

func (h *handler) login(ctx *gin.Context) {
	csrfToken := xsrftoken.Generate("supermegasecret", "", "")
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": h.Conf.AuthCodeURL(csrfToken),
	})
}

func (h *handler) oauth(ctx *gin.Context) {
	ok := xsrftoken.Valid(ctx.Query("state"), "supermegasecret", "", "")
	if !ok {
		ctx.JSON(401, gin.H{
			"code":    401,
			"message": "Invalid CSRF token",
		})
		return
	}
	token, er := h.Conf.Exchange(context.Background(), ctx.Query("code"))
	if er != nil {
		ctx.JSON(401, gin.H{
			"code":    401,
			"message": "Invalid code",
		})
		return
	}
	fmt.Println(token)
	ctx.Header("Authorization", token.AccessToken)
	ctx.Redirect(307, "http://localhost:8000/callback")
}
