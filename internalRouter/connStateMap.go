package internalrouter

import(
	upstream "sshGateway/sshListener"
)

var ActiveConns map[string]ConnEntry

type ConnEntry struct {
	WantedRes upstream.RequestedResource
	backendUp bool
	backendAddr string
}



func init() {
	ActiveConns = make(map[string]ConnEntry)

	//gotta move in controller hook-ins here for booting up the backends




}