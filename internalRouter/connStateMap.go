package internalrouter

import(
	frontend "sshGateway/sshListener"
)

var ActiveConns map[string]ConnEntry

type ConnEntry struct {
	WantedRes frontend.RequestedResource
	backendUp bool
	backendAddr string
}



func init() {
	ActiveConns = make(map[string]ConnEntry)

	//gotta move in controller hook-ins here for booting up the backends




}