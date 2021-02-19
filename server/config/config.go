package config

import (
	"database/sql"
	"fmt"
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
	DB  *sql.DB
}

// ProvideConfig creates a Config struct from environment variables and various
// other sources supported by Viper. Dynamic components, such as the DB, will
// then be instantiated.
func ProvideConfig() *Config {
	// Load env config
	viper.SetDefault("Port", 3000)
	viper.SetDefault("StaticDir", "./static")
	viper.SetDefault("DbPath", "data.db")
	viper.SetDefault("DbParams", "cache=shared&immutable=true&mode=ro&nolock=true")
	viper.AutomaticEnv()
	ec := EnvConfig{}
	err := viper.Unmarshal(&ec)
	if err != nil {
		log.Fatal(err)
	}

	// Instantiate dynamic components
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?%s", ec.DbPath, ec.DbParams))
	if err != nil {
		log.Fatal(err)
	}
	// Ensure DB is accessible
	_, err = db.Query("SELECT 1")
	if err != nil {
		log.Fatal(err)
	}
	c := Config{Env: &ec, DB: db}

	return &c
}
