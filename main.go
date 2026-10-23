package main

import (
	"sshGateway/internalRouter/adressResolver"
	"sshGateway/config"
	sshlistener "sshGateway/sshListener"
	//sshdownstream "sshGateway/downstream"
	kubImpl "sshGateway/internalRouter/kubernetesImpl"
)

func main() {
	config.InitConfig()
	println("Starting SSH listener")
	
	
	resolverKub := addressResolver.KubernetesResolver{}
	kc := &kubImpl.KubContext{}

	//have to set proper info here
	kc.InformerCache("kubeconfig", "")
	//resolverStatic := addressResolver.StaticResolver{}

	

	sshlistener.StartListener(9090)
	println("Starting downstream handshake")
	
	//sshdownstream.DownstreamHandshake()

}
