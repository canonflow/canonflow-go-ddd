package http

import (
	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	App            *gin.Engine
	Handler        *UserHandler
	AuthMiddleware *gin.HandlerFunc
}

func NewUserRoute(app *gin.Engine, handler *UserHandler, authMiddleware *gin.HandlerFunc) *UserRoute {
	return &UserRoute{
		App:            app,
		Handler:        handler,
		AuthMiddleware: authMiddleware,
	}
}

func (route *UserRoute) Init() {
	auth := route.App.Group("auth")
	{
		auth.POST("/signup", route.Handler.SignUp)
	}
}
