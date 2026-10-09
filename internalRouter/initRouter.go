package internalrouter

import(
	c "sshGateway/config"
)


func InitRouter() {
	var stat StaticResolver
	
	for _, route := range c.Gconfig.StaticRoutes {
		stat.AddStaticResolveEntry(route.Resource, route.Target)
	}

	if len(c.)
}