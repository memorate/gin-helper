package restful

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type Route struct {
	Method       string                 //Method is one of the following: GET,PUT,POST,DELETE. required
	Path         string                 //Path contains a path pattern. required
	ResourceFunc func(ctx *gin.Context) //the func this API calls.
}

type routerGroup struct {
	root     *routerGroup
	routers  []Route
	basePath string // the common prefix for all routes in this group
}

type Router interface {
	//URLPatterns returns route
	URLPatterns() []Route
}

func (r *Route) RegisterRoute(e *gin.Engine) {
	e.Handle(r.Method, r.Path, r.ResourceFunc)
}

func NewRouterGroup(basePath string) *routerGroup {
	rg := &routerGroup{basePath: basePath}
	rg.root = rg
	return rg
}

func (rg *routerGroup) GetRouters() []Route {
	return rg.routers
}

func (rg *routerGroup) Group(path string) *routerGroup {
	bp := rg.basePath + path
	if rg.basePath == "/" {
		bp = path
	}
	return &routerGroup{
		basePath: bp,
		root:     rg.root,
	}
}

func (rg *routerGroup) POST(relativePath string, handlers func(ctx *gin.Context)) *routerGroup {
	route := Route{
		Method:       http.MethodPost,
		Path:         rg.basePath + relativePath,
		ResourceFunc: handlers,
	}
	rg.combine(route)
	return rg
}

func (rg *routerGroup) GET(relativePath string, handlers func(ctx *gin.Context)) *routerGroup {
	route := Route{
		Method:       http.MethodGet,
		Path:         rg.basePath + relativePath,
		ResourceFunc: handlers,
	}
	rg.combine(route)
	return rg
}

func (rg *routerGroup) combine(route Route) {
	root := rg.root
	if rg.root == nil {
		root = rg
	}
	root.routers = append(root.routers, route)
}

func (rg *routerGroup) GetRouteList() []Route {
	return rg.routers
}

func GetUrlPatterns(schema interface{}) ([]Route, error) {
	v, ok := schema.(Router)
	if !ok {
		return []Route{}, fmt.Errorf("can not register APIs to server: %s", reflect.TypeOf(schema).String())
	}
	return v.URLPatterns(), nil
}
