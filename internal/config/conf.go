package config

import (
	"log"
	"sync"
)

type Config struct {
	HttpConf struct {
		Port string `yaml:"port" env-default:"8282"`
	}
	Postgres struct {
	}
	Redis struct {
	}
	Debug *bool `yaml:"debug"`
}

var appConfig *Config
var once sync.Once

func LoadConfig() *Config {
	once.Do(func() {
		// TODO get app logger
		log := log.Default()
		log.Println("read config..")
		appConfig = &Config{}
	})
	return appConfig
}
