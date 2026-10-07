package config

import(
	"github.com/go-yaml/yaml"
	"os"
)

type StaticRoutes struct {
	Source string `yaml:"endpoint"`
}

type GlobalConfig struct {
	StaticRoutes 		[]StaticRoutes 	`yaml:"staticRoutes"`
	ListenEndpoint 		  string 	`yaml:"listenEndpoint"`
	KubernetesControlEndp string 	`yaml:"kubernetesControlEndp"`
	ecdhsKey 			  string 	`yaml:"ecdhsKey"`
}

func InitConfig(){
	file, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}
	defer file.Close()



	decoder := yaml.NewDecoder(file)
	var config GlobalConfig
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}


}