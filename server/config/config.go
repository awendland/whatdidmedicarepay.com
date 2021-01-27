package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config contains all configuration parameters for the server
type Config struct {
	Port      int
	StaticDir string
	DbPath    string
}

// ProvideConfig instantiates a Config struct from environment variables and various
// other sources supported by Viper.
func ProvideConfig() *Config {
	viper.SetDefault("Port", 3000)
	viper.SetDefault("StaticDir", "./static")
	viper.AutomaticEnv()
	c := Config{}
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatal(err)
	}
	return &c
}
