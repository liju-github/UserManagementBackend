package config

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	DBUSER     string
	DBPASSWORD string
	DBHOST     string
	DBPORT     string
	DBNAME     string
}

func EnvConfig() Env {
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	viper.AutomaticEnv()

	var env Env

	env.DBUSER = viper.GetString("DBUSER")
	env.DBPASSWORD = viper.GetString("DBPASSWORD")
	env.DBHOST = viper.GetString("DBHOST")
	env.DBPORT = viper.GetString("DBPORT")
	env.DBNAME = viper.GetString("DBNAME")

	return env
}
