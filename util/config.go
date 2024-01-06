package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseSource        string        `mapstructure:"DB_SOURCE"`
	DatabaseDriver        string        `mapstructure:"DB_DRIVER"`
	ServerAddress         string        `mapstructure:"SERVER_ADDRESS"`
	TokenSecret           string        `mapstructure:"TOKEN_SECRET"`
	TokenDuration         time.Duration `mapstructure:"TOKEN_DURATION"`
	AWS_REGION            string        `mapstructure:"AWS_REGION"`
	AWS_ACCESS_KEY_ID     string        `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWS_SECRET_ACCESS_KEY string        `mapstructure:"AWS_SECRET_ACCESS_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
