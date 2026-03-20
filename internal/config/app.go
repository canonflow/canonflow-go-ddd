package config

import (
	"github.com/IBM/sarama"
	"github.com/canonflow/canonflow-go-ddd/internal/contract"
	userHttp "github.com/canonflow/canonflow-go-ddd/internal/domain/user/delivery/http"
	gateway "github.com/canonflow/canonflow-go-ddd/internal/domain/user/gateway/messaging"
	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/repository"
	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/usecase"
	"github.com/canonflow/canonflow-go-ddd/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *gin.Engine
	Redis    *redis.Client
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	Producer sarama.SyncProducer
}

func Bootstrap(config *BootstrapConfig) {
	// TODO: Setup all repositories
	userRepository := repository.NewUserRepositoryImpl(config.DB)

	// TODO: Setup all producers
	var userProducer contract.ProducerContract

	if config.Producer != nil {
		userProducer = gateway.NewUserProducer(config.Producer, config.Log)
	} else {
		config.Log.Warn("producer is not initialized")
	}

	// TODO: Setup all use cases
	userUsecase := usecase.NewUserUsecase(config.DB, config.Log, config.Config, userRepository, userProducer)

	// TODO: Setup all handler
	userHandler := userHttp.NewUserHandler(userUsecase, config.Config)

	// TODO: Setup all middlewares
	authMiddleware := middleware.AuthMiddleware(config.Config.GetString("JWT_SECRET_KEY"))

	// TODO: Setup Routes
	userRoute := userHttp.NewUserRoute(config.App, userHandler, &authMiddleware)
	userRoute.Init()
}
