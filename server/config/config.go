package config

import (
	"BaselineCheck/server/repository"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port int                  `mapstructure:"port" json:"port" yaml:"port"`
	Host string               `mapstructure:"host" json:"host" yaml:"host"`
	Repo *repository.DBConfig `mapstructure:"repo" json:"repo" yaml:"repo"`
}

func InitConfig(confFile string) (*Config, error) {
	var conf Config
	v := viper.New()
	v.SetConfigFile(confFile)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
		return nil, err
	}

	if err := v.Unmarshal(&conf); err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
		return nil, err

	}
	return &conf, nil
}
