package log

import (
	log "github.com/sirupsen/logrus"
	"os"
	"github.com/spf13/viper"
)

func init() {
	var env = viper.GetString("env")

	if env != "dev" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	log.SetOutput(os.Stdout)

	switch lvl := viper.GetString("log_level"); lvl {
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}

	log.Info("Log level is " , log.GetLevel())
}
