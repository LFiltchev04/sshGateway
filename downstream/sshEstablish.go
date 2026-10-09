package downstream

import (
	//"io"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
	//"net"
	//"os"
)


func  DownstreamHandshake() {
	config := &ssh.ClientConfig{
		User: "llf",
		Auth: []ssh.AuthMethod{
			ssh.Password("alabala1"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", "192.168.1.40:22", config)
	if err != nil {
		log.Fatalf("Failed to dial SSH server: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create SSH session: %v", err)
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	session.RequestPty("xterm", 24, 80, ssh.TerminalModes{})
	session.Shell()

	println("completed")
	
	session.Wait()

	
}