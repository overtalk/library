package gate

import "github.com/gin-gonic/gin"

// AddGroup : Add route group
func (gate *Gate) AddGroup(relativePath string, handlers ...gin.HandlerFunc) {
	gate.groups[relativePath] = gate.engine.Group(relativePath, handlers...)
}

// GET add get handler
func (gate *Gate) GET(relativePath string, handler interface{}, group ...string) {
	gate.routerGroup(group...).GET(relativePath, wrap(handler))
}

// POST add post handler
func (gate *Gate) POST(relativePath string, handler interface{}, group ...string) {
	gate.routerGroup(group...).POST(relativePath, wrap(handler))
}

// routerGroup get route group
func (gate *Gate) routerGroup(group ...string) *gin.RouterGroup {
	switch len(group) {
	case 1:
		return gate.groups[group[0]]
	default:
		return &gate.engine.RouterGroup
	}
}
