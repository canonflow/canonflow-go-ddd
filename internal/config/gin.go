package config

import (
	"fmt"
	"net/http"

	"github.com/canonflow/canonflow-go-ddd/pkg/response"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func PanicRecovery(log *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("%s \"%s\" - Panic Occured: %+v", ctx.Request.Method, ctx.Request.URL, err)

				//! Return a unified error response
				ctx.JSON(http.StatusInternalServerError, response.BaseErrorResponse{
					Code:   http.StatusInternalServerError,
					Status: "Internal Server Error",
					Error:  fmt.Sprintf("%v", err),
				})

				ctx.Abort()
			}
		}()

		ctx.Next()
	}
}

func NewGin(config *viper.Viper, log *logrus.Logger) *gin.Engine {
	app := gin.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	app.Use(PanicRecovery(log))

	app.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, response.BaseErrorResponse{
			Code:   http.StatusNotFound,
			Status: "Not Found",
			Error:  fmt.Sprintf("%s not found", ctx.Request.URL),
		})
	})

	return app
}
