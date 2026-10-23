package downstream

import (
	//"io"
	"io"
	"log"

	"golang.org/x/crypto/ssh"
	//"net"
	//"os"
)

func DownstreamHandshake(writer io.ReadWriter, reader io.Reader, ready chan<- *ssh.Session) {
	config := &ssh.ClientConfig{
		User: "llf",
		Auth: []ssh.AuthMethod{
			ssh.Password("alabala1"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", "192.168.1.40:22", config)
	if err != nil {
		log.Printf("Failed to dial SSH server: %v", err)
		ready <- nil
		return
	}

	session, err := client.NewSession()
	if err != nil {
		log.Printf("Failed to create SSH session: %v", err)
		client.Close()
		ready <- nil
		return
	}

	session.Stdin = reader
	session.Stdout = writer
	session.Stderr = writer

	if err := session.RequestPty("xterm-256color", 24, 80, ssh.TerminalModes{}); err != nil {
		log.Printf("Failed to request downstream PTY: %v", err)
		session.Close()
		client.Close()
		ready <- nil
		return
	}
	if err := session.Shell(); err != nil {
		log.Printf("Failed to start downstream shell: %v", err)
		session.Close()
		client.Close()
		ready <- nil
		return
	}

	println("completed")

	ready <- session
	if err := session.Wait(); err != nil {
		log.Printf("Downstream session ended: %v", err)
	}
	client.Close()
}
