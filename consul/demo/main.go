package main

import (
	"fmt"
	"log"
	"time"

	"web-layout/utils/consul"
)

func main() {
	go func() {
		time.Sleep(3 * time.Second)
		service2()
	}()

	go func() {
		time.Sleep(3 * time.Second)
		service1()
	}()

	listen()
}

func listen() {
	fmt.Println("开始监听服务 --- service1 & service2")
	consulAddr := "127.0.0.1:8500"
	r, err := consul.NewClient(consulAddr)
	if err != nil {
		log.Fatal(err)
	}

	r, err = r.ServiceDiscovery(&consul.DiscoveryConfig{
		ServerType: "service1",
		Tags:       []string{"0.98", "QQ"},
		Min:        1,
	}, &consul.DiscoveryConfig{
		ServerType: "service2",
		Tags:       []string{"0.98", "QQ"},
		Min:        1,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("开始等待服务上线")
	r.Wait()
	fmt.Println("所有依赖服务都上线")
	fmt.Println("接下来监听服务变化")

	queue := r.Watch()
	for {
		select {
		case data := <-queue:
			fmt.Println("服务状态发送变化了")
			fmt.Println(data)
		default:
			time.Sleep(time.Second)
		}
	}

}

func service1() {
	const (
		ServerType = "service1"
		ServiceID  = ServerType + "-1"
		checkPort  = 9000
		consulAddr = "127.0.0.1:8500"
	)

	service := &consul.RegistryConfig{
		IP:         "127.0.0.1",
		ID:         ServiceID,
		Port:       944,
		ServerType: ServerType,
		Tags:       []string{"0.98", "QQ"},
	}
	r, err := consul.NewClient(consulAddr)
	if err != nil {
		log.Fatal(err)
	}

	r.ServiceRegistry(checkPort, service)

	fmt.Println("开始注册服务 --- ", ServerType, " 20s 之后注销服务")
	if err := r.Register(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(20 * time.Second)

	fmt.Println("注销服务 --- ", ServerType)
	r.DeRegister()
}

func service2() {
	const (
		ServerType = "service2"
		ServiceID  = ServerType + "-1"
		checkPort  = 9001
		consulAddr = "127.0.0.1:8500"
	)

	service := &consul.RegistryConfig{
		IP:         "127.0.0.1",
		ID:         ServiceID,
		Port:       99992,
		ServerType: ServerType,
		Tags:       []string{"0.98", "QQ"},
	}
	r, err := consul.NewClient(consulAddr)
	if err != nil {
		log.Fatal(err)
	}

	r.ServiceRegistry(checkPort, service)

	fmt.Println("开始注册服务 --- ", ServerType, " 20s 之后注销服务")
	if err := r.Register(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(20 * time.Second)

	fmt.Println("注销服务 --- ", ServerType)
	r.DeRegister()
}
