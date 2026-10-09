package main

import (

	"sshGateway/config"
	sshlistener "sshGateway/sshListener"
	//sshdownstream "sshGateway/downstream"

)

func main() {
	config.InitConfig()
	println("Starting SSH listener")
	sshlistener.StartListener(9090)
	//println("Starting downstream handshake")
	
	
	
	//sshdownstream.DownstreamHandshake()

}
