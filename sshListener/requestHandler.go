package sshlistener

import (
	"bytes"

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

			var staticB [32]byte
			var buf bytes.Buffer
			buf.Grow(1024)
			for{
				
				n, err := ch.Read(staticB[:])
				if err != nil {
					println("Failed to read from shell channel: ", err)
					break
				}
				
				
				buf.Write(staticB[:n])
				
				
				ch.Write(buf.Bytes())
			}

	
	}
	
}