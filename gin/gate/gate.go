package gate

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Gate defines gateway
type Gate struct {
	srv    *http.Server
	port   int
	engine *gin.Engine
	groups map[string]*gin.RouterGroup
}

// NewGate is the constructor of Gate
func NewGate(port int, m ...gin.HandlerFunc) *Gate {
	engine := gin.New()
	engine.Use(m...)
	return &Gate{
		port:   port,
		engine: engine,
		groups: make(map[string]*gin.RouterGroup),
	}
}

// Start start gate server
func (gate *Gate) Start() {
	gate.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", gate.port),
		Handler: gate.engine,
	}

	fmt.Printf("Server Start on port : %d\n", gate.port)
	if err := gate.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic("Server Start Error : " + err.Error())
	}
}

// Shutdown stop the server with a timeout of 5 seconds
func (gate *Gate) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := gate.srv.Shutdown(ctx); err != nil {
		panic("Server Shutdown Error : " + err.Error())
	}
	fmt.Println("Server Exit")
}
