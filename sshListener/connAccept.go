package sshlistener

import(
	"golang.org/x/crypto/ssh"
	"sshGateway/downstream"
)

func HandleInputChannel(conn *ssh.ServerConn, reqs <-chan *ssh.Request, chans <-chan ssh.NewChannel) {

	for ch := range chans{
		nChan, rChan, err := ch.Accept()
		if err != nil {
			println("Failed to accept channel: ", err)
			return
		}
		nChan.Write([]byte("Hello from server"))

		select {
			case r := <-rChan:
				println("Request channel sent to handler")
				if r != nil {
					println("Received request")

					r.Reply(true, nil)
					
					downstream.DownstreamHandshake(nChan.Stderr())
				
				}
		}


		println("Channel handling completed")
	}
}