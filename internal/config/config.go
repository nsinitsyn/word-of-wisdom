package config

import (
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Addr       string `yaml:"addr" env-required:"true"`
	Complexity uint8  `yaml:"complexity" env-required:"true"`
}

type ClientConfig struct {
	ServerAddr string `yaml:"server-addr" env-required:"true"`
}

var serverConfig *ServerConfig
var clientConfig *ClientConfig
var once sync.Once

func GetServerConfig() *ServerConfig {
	once.Do(func() {
		path, ok := os.LookupEnv("CONFIG_PATH")
		if !ok {
			path = "server.yaml"
		}

		serverConfig = &ServerConfig{}
		if err := cleanenv.ReadConfig(path, serverConfig); err != nil {
			log.Fatal(err)
		}
	})
	return serverConfig
}

func GetClientConfig() *ClientConfig {
	once.Do(func() {
		path, ok := os.LookupEnv("CONFIG_PATH")
		if !ok {
			path = "client.yaml"
		}

		clientConfig = &ClientConfig{}
		if err := cleanenv.ReadConfig(path, clientConfig); err != nil {
			log.Fatal(err)
		}
	})
	return clientConfig
}
