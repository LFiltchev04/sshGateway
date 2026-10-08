package sshlistener

import(
	"golang.org/x/crypto/ssh"
)

func RequestHandle(req ssh.Request, ch ssh.Channel) {
	


	switch req.Type{
		case "pty-req":
			println("Handling pty-req")
			if req.WantReply {
				err := req.Reply(true, nil)
				if err != nil {
					println("Failed to reply to pty-req: ", err)
				}
				println("Replying to pty-req")
			}

		case "shell":
			println("Handling shell")
			if req.WantReply {
				err := req.Reply(true, nil)
				if err != nil {
					println("Failed to reply to shell: ", err)
				}
				println("Replying to shell")
			}

			for{
				buf := make([]byte, 1024)
				n, err := ch.Read(buf)
				if err != nil {
					println("Failed to read from shell channel: ", err)
					break
				} else {
					println("Read from shell channel: ", string(buf[:n]))
				}

				buf[n] = '\n'
				buf[n+1] = '\r'
				ch.Write(buf[:n+2])
			}

	
	}
	
}