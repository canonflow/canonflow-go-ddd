package config

import (
	"strings"
	"time"

	"github.com/canonflow/canonflow-go-ddd/internal/factory"
	"github.com/canonflow/canonflow-go-ddd/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var availableDriver = []string{"mysql", "postgre"}

func NewDatabase(config *viper.Viper, log *logrus.Logger) *gorm.DB {
	driver := strings.ToLower(config.GetString("DB_DRIVER"))
	username := config.GetString("DB_USER")
	password := config.GetString("DB_PASSWORD")
	host := config.GetString("DB_HOST")
	port := config.GetInt("DB_PORT")
	database := config.GetString("DB_NAME")
	idleConnection := config.GetInt("DB_IDLE")
	maxConnection := config.GetInt("DB_MAX")
	maxLifeTimeConnection := config.GetInt("DB_LIFETIME")

	if !utils.SliceContains(availableDriver, driver) {
		log.Fatal("Unsupported Database Driver")
	}

	databaseDriver := factory.NewDatabaseFactory(driver, username, password, host, port, database)
	if databaseDriver == nil {
		log.Fatal("Unsupported Database Driver")
	}

	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	dsn := databaseDriver.GetDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	return db
}

type logrusWriter struct {
	Logger *logrus.Logger
}

func (l *logrusWriter) Printf(message string, args ...interface{}) {
	l.Logger.Tracef(message, args...)
}
