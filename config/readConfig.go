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
type DbType struct {
	AccessPath string `yaml:"accessPath"`
	Type string `yaml:"type"`
}

type GlobalConfig struct {
	ListenEndpoint        string         `yaml:"listenEndpoint"`
	StaticRoutes          []StaticRoutes `yaml:"staticRoutes"`
	KubernetesControlEndp string         `yaml:"kubernetesControlEndp"`
	EcdsaKey              string         `yaml:"ecdsaKey"`
	DatabasePath          []DbType       `yaml:"databasePath"`
}


//global config struct as coming over from config yaml
var Gconfig GlobalConfig

func InitConfig() {
	file, err := os.Open("./config.yaml")
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


func (g *GlobalConfig) GetListenEndpoint() string {
	return g.ListenEndpoint
}

func (g* GlobalConfig) GetSSHKey() (string, error) {
	file, err := os.Open(g.EcdsaKey)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func isStaticRoute(resource string) bool {
	for _, route := range Gconfig.StaticRoutes {
		if route.Resource == resource {
			return true
		}
	}
	return false
}

func (g *GlobalConfig) GetStaticRouteTarget(resource string) (string, bool) {
	for _, route := range g.StaticRoutes {
		if route.Resource == resource {
			return route.Target, true
		}
	}
	return "", false
}

func (g *GlobalConfig) GetAllStaticRoutes() []StaticRoutes {
	return g.StaticRoutes
}