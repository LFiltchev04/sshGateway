package internalrouter

import(
	upstream "sshGateway/sshListener"
)





type ConnEntry struct {
	WantedRes upstream.RequestedResource
	backendUp bool
	backendAddr string
}