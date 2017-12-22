package config

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	Port       = "port"
	Env        = "env"
	DefaultEnv = "dev"
)

func init() {
	log.Info("Reading in config files...")

	pflag.String(Env, DefaultEnv, "environment config value to use")

	viper.BindPFlags(pflag.CommandLine)

	pflag.Parse()

	log.Info("Env is " + viper.GetString(Env))

	viper.SetConfigName("config." + viper.GetString(Env))
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../config")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}

	// Real-Time Config Changes
	// go func() {
	// 	for {
	// 		time.Sleep(time.Second * 5)
	// 		viper.WatchConfig()
	// 		viper.OnConfigChange(func(e fsnotify.Event) {
	// 			log.Println("config file changed: ", e.Name)
	// 		})
	// 	}
	// }()
}
