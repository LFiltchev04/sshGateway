package sshlistener

import (
	"log"

	d "sshGateway/downstream"

	"golang.org/x/crypto/ssh"
)

func HandleInputChannel(conn *ssh.ServerConn, reqs <-chan *ssh.Request, chans <-chan ssh.NewChannel) {
	for reqs != nil || chans != nil {
		select {

			case req, ok := <-reqs:
				if !ok {
					reqs = nil
					continue
				}
				if req != nil {
					req.Reply(true, nil)
				}
			
			
			case ch, ok := <-chans:
				if !ok {
					chans = nil
					continue
				}
			
				channel, channelRequests, err := ch.Accept()
				if err != nil {
					log.Printf("Failed to accept channel: %v", err)
					continue
				}

				go func(channel ssh.Channel, requests <-chan *ssh.Request) {
					defer channel.Close()
					//channel.Write([]byte("\n\n===>          THIS IS BEING FORWARDED         <===\n\n"))

					sessionReady := make(chan *ssh.Session, 1)
					go d.DownstreamHandshake(channel.Stderr(), channel, sessionReady)
					
					downstream := <-sessionReady
					if downstream == nil {
						return
					}

					for request := range requests {
						if request == nil {
							continue
						}

						switch request.Type {
							case "pty-req":
								ok, err := downstream.SendRequest("pty-req", request.WantReply, request.Payload)
								if err != nil {
									log.Printf("Failed to forward pty-req: %v", err)
								}
							
								if request.WantReply {
									request.Reply(ok, request.Payload)
								}
							
							case "shell":

								request.Reply(true, nil)
								ok, err := downstream.SendRequest("shell", request.WantReply, nil)
								if err != nil {
									log.Printf("Failed to forward shell: %v", err)
								}

								if request.WantReply {
									request.Reply(ok, nil)
								}
					
							case "window-change":
								ok, err := downstream.SendRequest("window-change", request.WantReply, request.Payload)
								if err != nil {
									log.Printf("Failed to forward window change: %v", err)
								}

								if request.WantReply {
									request.Reply(ok, nil)
								}



							case "exec":
								ok, err := downstream.SendRequest("exec", request.WantReply, request.Payload)
								if err != nil {
									log.Printf("Failed to forward exec: %v", err)
								}

								if request.WantReply {
									request.Reply(ok, nil)
								}

							default:
								resp, err := downstream.SendRequest(request.Type, request.WantReply, nil)
								if err != nil {
									log.Printf("Failed to forward request: %v", err)
								}
						
								if request.WantReply == true{
									request.Reply(resp, nil)
								}
							
							}
					}
			}(channel, channelRequests)
		}
	}
}
