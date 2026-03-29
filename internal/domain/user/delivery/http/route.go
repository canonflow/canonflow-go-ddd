package http

import (
	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	App                   *gin.Engine
	Handler               *UserHandler
	AuthMiddleware        *gin.HandlerFunc
	RateLimiterMiddleware *gin.HandlerFunc
}

func NewUserRoute(app *gin.Engine, handler *UserHandler, authMiddleware *gin.HandlerFunc, rateLimiterMiddleware *gin.HandlerFunc) *UserRoute {
	return &UserRoute{
		App:                   app,
		Handler:               handler,
		AuthMiddleware:        authMiddleware,
		RateLimiterMiddleware: rateLimiterMiddleware,
	}
}

func (route *UserRoute) Init() {
	auth := route.App.Group("api/v1/auth")
	{
		auth.POST("/signup", route.Handler.SignUp)
		auth.POST("/signin", route.Handler.SignIn)

		authMiddleware := auth.Group("")
		{
			authMiddleware.Use(*route.AuthMiddleware)
			rateLimiterMiddleware := authMiddleware.Group("")
			{
				rateLimiterMiddleware.Use(*route.RateLimiterMiddleware)
				rateLimiterMiddleware.GET("/me", route.Handler.Me)
			}
			authMiddleware.POST("/signout", route.Handler.SignOut)
		}
	}
}
