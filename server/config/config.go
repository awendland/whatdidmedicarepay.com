package config

import (
	"log"

	"github.com/spf13/viper"
)

// EnvConfig contains all configuration parameters for the server that can
// be provided by the environment.
type EnvConfig struct {
	Port      int
	StaticDir string
	DbPath    string
	DbParams  string
}

// Config contains dynamic configuration items (like the DB pool) as well as
// provided configuration details from the environment.
type Config struct {
	Env *EnvConfig
}

// ProvideConfig creates a Config struct from environment variables and various
// other sources supported by Viper. Dynamic components, such as the DB, will
// then be instantiated.
func ProvideConfig() *Config {
	// Load env config
	viper.SetDefault("Port", 3000)
	viper.SetDefault("StaticDir", "./static")
	viper.AutomaticEnv()
	ec := EnvConfig{}
	err := viper.Unmarshal(&ec)
	if err != nil {
		log.Fatal(err)
	}
	return &c
}
