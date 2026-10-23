package internalrouter

import R "sshGateway/internalRouter/adressResolver"

type RequestedResource struct{
	Uname string
	EnvName string
	EnvVersion string
	Resolver R.DownstreamResolve
}


type ConnEntry struct {
	WantedRes RequestedResource
	backendUp bool
	backendAddr string
}

var ActiveConns map[string]ConnEntry
