package addressResolver

import (
	ir "sshGateway/internalRouter"
	"log"
)

type DownstreamResolve interface {
	Resolve(req ir.RequestedResource) string
	StartResource() error
}

type KubernetesResolver struct{}
func (k KubernetesResolver) Resolve(req ir.RequestedResource) string {
	//dead certain this will be somehow wrong
	fqdn := req.EnvName + req.EnvVersion + "." + req.Uname + ".svc.cluster.local"
	return fqdn
}

func (k KubernetesResolver) StartResource(req *ir.RequestedResource) error {
	log.Printf("Starting resource for %s", req.EnvName)
	

	return nil
}

type StaticResolver struct{
	allowed map[string]string
}
func (s StaticResolver) Resolve(req ir.RequestedResource) string {
	return req.EnvName + "." + req.Uname
}

func (s *StaticResolver) AddStaticResolveEntry(key, value string) {
	s.allowed[key] = value
}

func (s *StaticResolver) RemoveStaticResolveEntry(key string) {
	delete(s.allowed, key)
}




type dbImpl interface {
	isPresent(key string, envName string) bool
	getRef(key string, envName string) string
}
