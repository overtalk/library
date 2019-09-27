package gate

import "github.com/gin-gonic/contrib/static"

// AddGroup : Add route group
func (gate *Gate) Static(urlPrefix, root string) {
	gate.engine.Use(static.Serve(urlPrefix, static.LocalFile(root, true)))
}
