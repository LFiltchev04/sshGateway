package internalrouter


type RequestedResource struct{
	Uname string
	EnvName string
	EnvVersion string
}


type ConnEntry struct {
	WantedRes RequestedResource
	backendUp bool
	backendAddr string
}

var ActiveConns map[string]ConnEntry
