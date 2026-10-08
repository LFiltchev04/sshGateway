package sshlistener

import(
	"golang.org/x/crypto/ssh"
)

func HandleInputChannel(conn *ssh.ServerConn, reqs <-chan *ssh.Request, chans <-chan ssh.NewChannel) {

	for ch := range chans{
		nChan, rChan, err := ch.Accept()
		if err != nil {
			println("Failed to accept channel: ", err)
			return
		}
		nChan.Write([]byte("Hello from server"))

		func (rChan <-chan *ssh.Request) {
			for req := range rChan {
				go RequestHandle(*req, nChan)
			}
		}(rChan)
	}
}