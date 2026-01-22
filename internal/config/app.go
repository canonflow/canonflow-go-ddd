package config

import (
	userHttp "github.com/canonflow/canonflow-go-ddd/internal/domain/user/delivery/http"
	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/repository"
	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/usecase"
	"github.com/canonflow/canonflow-go-ddd/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *gin.Engine
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// TODO: Setup all repositories
	userRepository := repository.NewUserRepositoryImpl(config.DB)

	// TODO: Setup all use cases
	userUsecase := usecase.NewUserUsecase(config.DB, config.Log, config.Config, userRepository)

	// TODO: Setup all handler
	userHandler := userHttp.NewUserHandler(userUsecase)

	// TODO: Setup all middlewares
	authMiddleware := middleware.AuthMiddleware(config.Config.GetString("JWT_SECRET_KEY"))

	// TODO: Setup Routes
	userRoute := userHttp.NewUserRoute(config.App, userHandler, &authMiddleware)
	userRoute.Init()
}
