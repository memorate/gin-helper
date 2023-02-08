package gin_helper

import (
	"gin-helper/restful"
	"github.com/gin-gonic/gin"
	"sync"
)

var ginInstance = &GinInstance{}

type Schema struct {
	serverName string
	schema     interface{} // business application ptr instance
}

type GinInstance struct {
	version string
	schemas []*Schema
	mu      sync.Mutex
}

func RegisterSchema(severName string, schema interface{}) {
	ginInstance.registerSchema(severName, schema)
}

func (g *GinInstance) registerSchema(serverName string, structPtr interface{}) {
	schema := &Schema{
		serverName: serverName,
		schema:     structPtr,
	}
	g.mu.Lock()
	g.schemas = append(g.schemas, schema)
	g.mu.Unlock()
}

func RegisterUrlPatterns(engine *gin.Engine) error {
	for _, schema := range ginInstance.schemas {
		if routes, e := restful.GetUrlPatterns(schema.schema); e != nil {
			return e
		} else {
			for _, route := range routes {
				engine.Handle(route.Method, route.Path, route.ResourceFunc)
			}
		}
	}
	return nil
}
