package consul

import (
	"fmt"

	consulAPI "github.com/hashicorp/consul/api"
	consulWatch "github.com/hashicorp/consul/api/watch"
)

// ClientInterface defines the Client of Consul
// Registration/Registration service to Consul
// Listening to the service in Consul
type ClientInterface interface {
	Wait()                          // Wait for a specific service to come online
	Register() error                // Registration service to Consul
	DeRegister() error              // DeRegister service to Consul
	Watch() <-chan AvailableServers // Listening to the service in Consul
}

// AvailableServers defines available online services
type AvailableServers struct {
	ServerType string
	Servers    []string
}

// RegistryConfig is service registry config
type RegistryConfig struct {
	ID         string   // service id
	IP         string   // service addr
	Port       int      // service port
	ServerType string   // service type
	Tags       []string // service Tags
}

// DiscoveryConfig is service discovery config
type DiscoveryConfig struct {
	ServerType string   // target service type
	Tags       []string // target service tags
	Min        int      // minimum of available in wait
	// others
	watchChan chan AvailableServers
	plan      *consulWatch.Plan
}

func (discoveryConfig *DiscoveryConfig) handler(index uint64, raw interface{}) {
	if raw == nil {
		return
	}
	if entries, ok := raw.([]*consulAPI.ServiceEntry); ok {
		var servers []string
		for _, entry := range entries {
			// healthy check fail, continue anyway
			if entry.Checks.AggregatedStatus() != consulAPI.HealthPassing {
				continue
			}
			servers = append(servers, fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port))
		}
		discoveryConfig.watchChan <- AvailableServers{
			ServerType: discoveryConfig.ServerType,
			Servers:    servers,
		}
	}
}
