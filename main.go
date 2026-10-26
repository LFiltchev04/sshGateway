package main

import (
	"sshGateway/config"
	sshlistener "sshGateway/sshListener"
	//sshdownstream "sshGateway/downstream"
	kubImpl "sshGateway/internalRouter/kubernetesImpl"
)

func main() {
	config.InitConfig()
	println("Starting SSH listener")
	
	
	kc := &kubImpl.KubContext{}

	//have to set proper info here
	kc.InformerCache("kubeconfig", "")
	//resolverStatic := addressResolver.StaticResolver{}

	

	sshlistener.StartListener(9090, kc)
	println("Starting downstream handshake")
	
	//sshdownstream.DownstreamHandshake()

}
