package main

import (
	"log"
	"sshGateway/config"
	sshlistener "sshGateway/sshListener"
	//sshdownstream "sshGateway/downstream"
	kubImpl "sshGateway/internalRouter/kubernetesImpl"
)

func main() {
	config.InitConfig()
	println("Starting SSH listener")

	kc := &kubImpl.KubContext{}

	unblockChan := make(chan bool, 1)
	//have to set proper info here
	go kc.InformerCache("kubeconfig", "", unblockChan)
	//resolverStatic := addressResolver.StaticResolver{}

	//just waits until init is done
	if ready := <-unblockChan; !ready {
		log.Fatal("Kubernetes informer cache failed to sync")
	}

	sshlistener.StartListener(9090, kc)
	println("Starting downstream handshake")

	//sshdownstream.DownstreamHandshake()

}
