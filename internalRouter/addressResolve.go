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




type dbImpl interface {
	isPresent(key string) bool
	getRef(key string) string
}

type PostgresBackedResolver struct{
	typeName string
	resolvImpl dbImpl
}

func (d PostgresBackedResolver) Resolve(req frontend.RequestedResource) string {
	//not today
	if d.resolvImpl.isPresent(req.EnvName) {
		return d.resolvImpl.getRef(req.EnvName)
	}
	
	return "buzz off"
}
func (d *PostgresBackedResolver) SetDatabaseImpl(impl dbImpl) {
	d.resolvImpl = impl
}
