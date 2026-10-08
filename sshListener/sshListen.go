package sshlistener


import (
	//"crypto/ed25519"
	//"crypto/rand"
	//"encoding/pem"
	"fmt"
	"log"
	"net"
 
	"golang.org/x/crypto/ssh"

	"sshGateway/config"
)


func StartListener(port int) {
	
	key, err := ssh.ParsePrivateKey([]byte(config.Gconfig.EcdsaKey))
	if(err != nil){
		log.Fatal("Failed to parse private key: ", err)
		panic(err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", port, err)
	}
	defer listener.Close()

	log.Printf("Listening on port %d...", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go func() {
			var conf ssh.ServerConfig
			conf.NoClientAuth = false
			conf.AddHostKey(key)
			
			_, chans, reqs, err := ssh.NewServerConn(conn, &conf)
			if err != nil {
				log.Printf("Failed to handshake: %v", err)
				return
			}

			for rq := range reqs {
				log.Printf("Request type: %v", rq.Type)
				
				log.Printf("Request: %v", rq)
			}
			go ssh.DiscardRequests(reqs)

			for newChannel := range chans {
				newChannel.Reject(ssh.Prohibited, "No channels are allowed")
			}
		}()
	}
}