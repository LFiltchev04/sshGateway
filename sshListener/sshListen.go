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
	
	keyContent, err := config.Gconfig.GetSSHKey()
	if err != nil {
		log.Fatal("Failed to get SSH key: ", err)
		panic(err)
	}
	key, err := ssh.ParsePrivateKey([]byte(keyContent))
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

		go func(cn net.Conn) {
			var conf ssh.ServerConfig
			conf.PasswordCallback = AuthCb
			conf.NoClientAuth = false
			conf.AddHostKey(key)
			
			sshConn, chans, _, err := ssh.NewServerConn(cn, &conf)
			if err != nil {
				log.Printf("Failed to handshake: %v", err)
				return
			}

			
			go HandleInputChannel(sshConn, nil, chans)
			
			println("SSH connection established")


			

		}(conn)
	}
}