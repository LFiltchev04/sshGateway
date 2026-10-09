package internalrouter

import(
	"log/slog"
	c "sshGateway/config"
)

var supportedMap map[string]struct{} = make(map[string]struct{})

func InitRouter() {
	//janky. have to rework 
	supportedMap["postgres"] = struct{}{}
	var stat StaticResolver
	
	for _, route := range c.Gconfig.StaticRoutes {
		stat.AddStaticResolveEntry(route.Resource, route.Target)
	}

	if len(c.Gconfig.DatabasePath) > 0 {
		
		for _, db := range c.Gconfig.DatabasePath {
			
			if supportedMap[db.Type] != struct{}{} {
				
				slog.Warn("unsupported database type", "type", db.Type)
				continue
			
			}else{
		
				if db.AccessPath == "" {
					slog.Warn("database path is empty", "type", db.Type)
				}

				

			}
		}
	}








}