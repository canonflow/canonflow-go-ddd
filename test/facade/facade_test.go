package facade_test

import (
	"context"
	"os"
	"testing"

	"github.com/canonflow/canonflow-go-ddd/internal/config"
	"github.com/canonflow/canonflow-go-ddd/pkg/facade"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	//* Set working directory to project root once for all tests
	os.Chdir("../../")
	os.Exit(m.Run())
}

func TestLoggerFromFacade(t *testing.T) {
	viperConfig := config.NewViper()
	logrus := config.NewLogrus(viperConfig)

	facade.InitializeFacadeDependencies(viperConfig, logrus)

	ctx := context.Background()
	facadeContext := facade.BuildContext(ctx, "abiwgiae-9283sa7d6s-12juhgu3ye7yd")

	facadeContext.GetLogger().Info("This is an info log with request ID")

	assert.Equal(t, "1", "1")
}

func TestConfigFromFacade(t *testing.T) {
	viperConfig := config.NewViper()
	logrus := config.NewLogrus(viperConfig)

	facade.InitializeFacadeDependencies(viperConfig, logrus)

	ctx := context.Background()
	facadeContext := facade.BuildContext(ctx, "abiwgiae-9283sa7d6s-12juhgu3ye7yd")

	assert.NotNil(t, facadeContext.GetConfig())
	assert.Equal(t, viperConfig.GetString("WEB_PORT"), facadeContext.GetConfig().GetString("WEB_PORT"))
}

func TestLoggerFromContext(t *testing.T) {
	viperConfig := config.NewViper()
	logrus := config.NewLogrus(viperConfig)

	facade.InitializeFacadeDependencies(viperConfig, logrus)

	ctx := context.Background()
	facadeContext := facade.BuildContext(ctx, "abiwgiae-9283sa7d6s-12juhgu3ye7yd")

	logrusContextLogger := facade.GetLoggerFromContext(facadeContext.Context)

	assert.NotNil(t, logrusContextLogger)
	logrusContextLogger.Info("This is an info log with request ID from Context")
}

func TestConfigFromContext(t *testing.T) {
	viperConfig := config.NewViper()
	logrus := config.NewLogrus(viperConfig)

	facade.InitializeFacadeDependencies(viperConfig, logrus)

	ctx := context.Background()
	facadeContext := facade.BuildContext(ctx, "abiwgiae-9283sa7d6s-12juhgu3ye7yd")

	configFromContext := facade.GetConfigFromContext(facadeContext.Context)

	assert.NotNil(t, configFromContext)
	assert.Equal(t, viperConfig.GetString("WEB_PORT"), configFromContext.GetString("WEB_PORT"))
}
