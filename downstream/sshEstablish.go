package downstream


import(
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"log"
	"net"
	"os"
)


func  DownstreamHandshake() {
	socketPath := "192.168.1.40"
	if socketPath == "" {
		log.Fatal("SSH_AUTH_SOCK environment variable is not set")
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Fatalf("Failed to connect to SSH agent socket: %v", err)
	}
	defer conn.Close()

	agentClient := agent.NewClient(conn)
	
	config := &ssh.ClientConfig{
		User: "your-username",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agentClient.Signers),
		},
		// WARNING: For production, use a proper ssh.HostKeyCallback instead of InsecureIgnoreHostKey
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", "example.com:22", config)
	if err != nil {
		log.Fatalf("Failed to dial SSH server: %v", err)
	}
	defer client.Close()

	clientSession, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create SSH session: %v", err)
	}
	defer clientSession.Close()

	clientSession.Stdin = os.Stdin
	clientSession.Stdout = os.Stdout
	clientSession.Stderr = os.Stderr

	err = clientSession.RequestPty("xterm", 80, 40, ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	})
	if err != nil {
		log.Fatalf("Failed to request PTY: %v", err)
	}
	
	err = clientSession.Shell()
	if err != nil {
		log.Fatalf("Failed to start shell: %v", err)
	}
}