package internalrouter

import (
	frontend "sshGateway/sshListener"
)

type downstreamResolve interface {
	resolve(req frontend.RequestedResource) string
}

type KubernetesResolver struct{}
func (k KubernetesResolver) resolve(req frontend.RequestedResource) string {
	//dead certain this will be somehow wrong
	fqdn := req.EnvName + req.EnvVersion + "." + req.Uname + ".svc.cluster.local"
	return fqdn
}

type StaticResolver struct{
	allowed map[string]string 
}
func (s StaticResolver) resolve(req frontend.RequestedResource) string {
	return req.EnvName + "." + req.Uname
}

func (s *StaticResolver) AddStaticResolveEntry(key, value string) {
	s.allowed[key] = value
}

func (s *StaticResolver) RemoveStaticResolveEntry(key string) {
	delete(s.allowed, key)
}



type DatabaseBackedResolver struct{
	dbUrl string
}

func (d DatabaseBackedResolver) resolve(req frontend.RequestedResource) string {
	//not today

	return "buzz off"
}