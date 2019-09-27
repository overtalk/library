package consul

import (
	"fmt"
	"net/http"
	"sync"

	consulAPI "github.com/hashicorp/consul/api"
	consulWatch "github.com/hashicorp/consul/api/watch"
)

// Client defines the consul client
type Client struct {
	consulAddr string // consul address（127.0.0.1:8500）

	// service registration related
	registryConfig *RegistryConfig
	checkPort      int // service registration check port
	checkServer    *http.Server
	consulClient   *consulAPI.Client // consul Client

	// service discovery related
	once             sync.Once
	discoveryConfigs []*DiscoveryConfig
	watchChan        chan AvailableServers
}

// NewClient is the constructor of consul Client
func NewClient(consulAddr string) (*Client, error) {
	// service registry
	c, err := consulAPI.NewClient(&consulAPI.Config{Address: consulAddr})
	if err != nil {
		return nil, err
	}

	return &Client{
		consulAddr:   consulAddr,
		consulClient: c,
	}, nil
}

func (client *Client) ServiceRegistry(checkPort int, registryConfig *RegistryConfig) *Client {
	// construct check sever
	mux := http.NewServeMux()
	mux.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	checkServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", checkPort),
		Handler: mux,
	}

	client.checkPort = checkPort
	client.checkServer = checkServer
	client.registryConfig = registryConfig
	return client
}

func (client *Client) ServiceDiscovery(discoveryConfigs ...*DiscoveryConfig) (*Client, error) {
	// service discovery channel
	watchChan := make(chan AvailableServers, 100)

	// service discovery
	for _, sdConfig := range discoveryConfigs {
		// build plan
		params := make(map[string]interface{})
		params["type"] = "service"
		params["service"] = sdConfig.ServerType
		params["tag"] = sdConfig.Tags
		plan, err := consulWatch.Parse(params)
		if err != nil {
			return nil, err
		}
		plan.Handler = sdConfig.handler

		// bind plan to DiscoveryConfig
		sdConfig.watchChan = watchChan
		sdConfig.plan = plan
	}

	client.discoveryConfigs = discoveryConfigs
	client.watchChan = watchChan
	client.once = sync.Once{}
	return client, nil
}
