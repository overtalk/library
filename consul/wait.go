package consul

type serviceDetail struct {
	min     int
	current int
}

func (serviceDetail *serviceDetail) available() bool {
	return serviceDetail.current >= serviceDetail.min
}

// Wait wait for required service
func (client *Client) Wait() {
	// new services count
	services := make(map[string]*serviceDetail)
	for _, v := range client.discoveryConfigs {
		services[v.ServerType] = &serviceDetail{
			min: v.Min,
		}
	}

	if allAvailable(services) {
		return
	}

	for msg := range client.Watch() {
		detail := services[msg.ServerType]
		detail.current = len(msg.Servers)
		if allAvailable(services) {
			return
		}
	}
}

// check is all type services available
func allAvailable(servers map[string]*serviceDetail) bool {
	for _, server := range servers {
		if !server.available() {
			return false
		}
	}
	return true
}
