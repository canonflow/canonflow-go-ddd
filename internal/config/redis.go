package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewRedis(config *viper.Viper, log *logrus.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.GetString("REDIS_HOST"), config.GetInt("REDIS_PORT")),
		Password: config.GetString("REDIS_PASSWORD"),
		DB:       0,
		Protocol: 2,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}

	return rdb
}
