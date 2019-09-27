package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	. "web-layout/utils/gin/gate"
)

func main() {
	g := NewGate(8080)
	g.POST("/test", test)

	g.AddGroup("/group1")
	g.GET("/test1", test1, "/group1")
	g.GET("/test2", test2, "/group1")

	g.AddGroup("/group2")
	g.GET("/test1", test3, "/group2")
	g.GET("/test2", test4, "/group2")

	g.Start()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	g.Shutdown()
}

func test(ctx context.Context, req *map[string]interface{}) map[string]interface{} {
	fmt.Println(req)
	return map[string]interface{}{"func": "test"}
}
func test1(ctx context.Context) map[string]interface{} { return map[string]interface{}{"func": "test1"} }
func test2(ctx context.Context) map[string]interface{} { return map[string]interface{}{"func": "test2"} }
func test3(ctx context.Context) map[string]interface{} { return map[string]interface{}{"func": "test3"} }
func test4(ctx context.Context) map[string]interface{} { return map[string]interface{}{"func": "test4"} }
