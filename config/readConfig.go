package config

import (
	"github.com/go-yaml/yaml"
	"io"
	"os"
)

type StaticRoutes struct {
	Resource string `yaml:"resource"`
	Target string `yaml:"target"`
}

type GlobalConfig struct {
	ListenEndpoint        string         `yaml:"listenEndpoint"`
	StaticRoutes          []StaticRoutes `yaml:"staticRoutes"`
	KubernetesControlEndp string         `yaml:"kubernetesControlEndp"`
	EcdsaKey              string         `yaml:"ecdsaKey"`
}

var Gconfig GlobalConfig
func InitConfig() {
	file, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var config GlobalConfig
	err = decoder.Decode(&config)
	if err == io.EOF {
		return
	}
	if err != nil {
		panic(err)
	}

	Gconfig = config
}
