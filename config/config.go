package config

import (
	"time"

	"github.com/spf13/viper"
)

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Config struct {
	Mode                string        `mapstructure:"mode"`
	Redis               RedisConfig   `mapstructure:"redis"`
	KeyLength           int           `mapstructure:"key_length"`
	KeyRotationInterval time.Duration `mapstructure:"key_rotation_interval"`
	Port                int           `mapstructure:"port"`
	Host                string        `mapstructure:"host"`
}

var GlobalConfig Config

func LoadConfig(path string) error {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml") // Or "toml"

	viper.AutomaticEnv() // Read environment variables

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&GlobalConfig)
	return err
}
