package config

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type CanonflowLogrusFormatter struct{}

func (f *CanonflowLogrusFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b bytes.Buffer

	timestamp := entry.Time.Format(time.DateTime)
	level := strings.ToUpper(entry.Level.String())
	msg := entry.Message
	requestID := "-"

	if rid, ok := entry.Data["request_id"]; ok {
		requestID = fmt.Sprintf("%v", rid)
	}

	fmt.Fprintf(&b, "%s[%s] [request_id=%s] %s\n", level, timestamp, requestID, msg)

	return b.Bytes(), nil
}

func NewLogrus(config *viper.Viper) *logrus.Logger {
	log := &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &CanonflowLogrusFormatter{},
		// Formatter: &easy.Formatter{
		// 	TimestampFormat: "2006-01-02 15:04:05",
		// 	LogFormat:       "%lvl%[%time%] %msg%\n\n",
		// },
	}

	log.SetLevel(logrus.Level(config.GetInt32("LOG_LEVEL")))

	return log
}
