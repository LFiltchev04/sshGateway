package main

import (

	"sshGateway/config"
	sshlistener "sshGateway/sshListener"

)

func main() {
	config.InitConfig()
	sshlistener.StartListener(9090)

}
