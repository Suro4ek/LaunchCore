package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type middleware struct {
	conf *oauth2.Config
}

func NewOAuthMiddleware() *middleware {
	return &middleware{}
}

type AuthUser struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type authHeader struct {
	IDToken string `header:"Authorization"`
}

func (m *middleware) CheckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h1 := authHeader{}
		if er := ctx.ShouldBindHeader(&h1); er != nil {
			ctx.AbortWithStatusJSON(401, gin.H{
				"code":    401,
				"message": "Invalid token",
			})
			return
		}
		idTokenHeader := strings.Split(h1.IDToken, "Bearer ")
		if len(idTokenHeader) < 2 {
			ctx.AbortWithStatusJSON(401, gin.H{
				"code":    401,
				"message": "Invalid token",
			})
			return
		}
		idToken := idTokenHeader[1]
		if idToken == "" {
			ctx.AbortWithStatusJSON(401, gin.H{
				"code":    401,
				"message": "Invalid token",
			})
			return
		}
		client := github.NewClient(m.conf.Client(context.TODO(), &oauth2.Token{AccessToken: idToken}))
		user, _, err := client.Users.Get(oauth2.NoContext, "")
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("Failed to get user: %v", err))
			return
		}
		authUser := AuthUser{
			Login: *user.Login,
			Name:  *user.Name,
			Email: *user.Email,
		}
		ctx.Set("user", authUser)
		ctx.Next()
	}
}
