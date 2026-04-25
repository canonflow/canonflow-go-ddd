package facade

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type (
	Logger struct{}
	Config struct{}
)

var (
	config *viper.Viper
	logger *logrus.Logger
)

type Facade struct {
	Context context.Context
}

func InitializeFacadeDependencies(viper *viper.Viper, logrus *logrus.Logger) {
	config = viper
	logger = logrus
}

func BuildContext(ctx context.Context, requestId string) *Facade {
	ctx = context.WithValue(ctx, Config{}, config)
	logger := logger.WithField("request_id", requestId)
	ctx = context.WithValue(ctx, Logger{}, logger)

	fa := &Facade{
		Context: ctx,
	}

	return fa
}

func (f *Facade) GetConfig() *viper.Viper {
	if f.Context.Value(Config{}) == nil {
		return nil
	}
	return f.Context.Value(Config{}).(*viper.Viper)
}

func (f *Facade) GetLogger() *logrus.Entry {
	if f.Context.Value(Logger{}) == nil {
		return nil
	}
	return f.Context.Value(Logger{}).(*logrus.Entry)
}

func GetConfigFromContext(ctx context.Context) *viper.Viper {
	if ctx.Value(Config{}) == nil {
		return nil
	}
	return ctx.Value(Config{}).(*viper.Viper)
}

func GetLoggerFromContext(ctx context.Context) *logrus.Entry {
	if ctx.Value(Logger{}) == nil {
		return nil
	}
	return ctx.Value(Logger{}).(*logrus.Entry)
}
