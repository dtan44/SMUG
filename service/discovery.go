package service

const ()

var ()

//DiscoveryInterface defines service methods
type DiscoveryInterface interface {
	List() []string
}

//DiscoveryService defines registration service struct
type DiscoveryService struct {
}

//List show all services avaliable
func (ds DiscoveryService) List() []string {
	keys := make([]string, 0, len(serviceMap))
	for k := range serviceMap {
		keys = append(keys, k)
	}
	return keys
}
