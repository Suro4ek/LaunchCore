package config

import (
	"sync"

	"LaunchCore/pkg/logging"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	MySQL     MySQL     `yaml:"mysql"`
	Minecraft Minecraft `yaml:"minecraft"`
	OAuth2    OAuth2    `yaml:"oauth2"`
}

type MySQL struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"pass"`
	DB       string `yaml:"db"`
}

type Minecraft struct {
	Ports  string `yaml:"ports"`
	Memory string `yaml:"memory"`
	Type   string `yaml:"type" default:"PAPER"`
}

type OAuth2 struct {
	ClientSecret string `yaml:"client_secret"`
	ClientID     string `yaml:"client_id"`
	RedirectURL  string `yaml:"redirect_url"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
