package downstream

import (
	//"io"
	"io"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
	//"net"
	//"os"

	rsl "sshGateway/internalRouter/adressResolver"
)


//is not idempotent, has to be ran once 
func DownstreamHandshake(writer io.ReadWriter, reader io.Reader, req rsl.RequestedResource,ready chan<- *ssh.Session) { 
	config := &ssh.ClientConfig{
		User: "llf",
		Auth: []ssh.AuthMethod{
			ssh.Password("alabala1"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}


	var ctx rsl.ClientCtx

	//init happens here, add later

	
	//failed to resolve at a basic level, drops connection
	if !req.Resolver.Resolve(req){
		writer.Write([]byte("\n\n         resolve failed\n\n"))
		ready <- nil
		return
	}	

	err := req.Resolver.StartResource(req)
	if err != nil {
		log.Printf("start failed with err: %s", err.Error())
		ready <- nil
		return
	}


	//this is going to be the resolution mechanism 
	//it talks to downstreams ervices and it either times out or allows it to proceed while handling the timeout
	for{
		go func(resolutionChannel *chan bool) {
			// send response to timer 
			time.Sleep(15 * time.Second)
			println("timer ended ")
			if resolutionChannel != nil {
				*resolutionChannel <- true
			}

			//is ignored if there is no channel because it has completed, just exit quietly
		}(ctx.StatusChan)

		// registers the resolution channel with the k8s controller



		// waits on resolution
		result := <-*ctx.StatusChan
		if result {
			break //permits to continue downwards to ssh
		}else{ //it has to kill the downstream channel and print a message
			log.Printf("Resolution failed")

			ready <- nil
			return
		}
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
