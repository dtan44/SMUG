package service

import (
	"errors"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	servicePath = "/service/"
)

var ()

//DiscoveryInterface defines service methods
type DiscoveryInterface interface {
	List() []string
	Route(*http.Request) (*http.Response, []byte, error)
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

//Route sends request to service
func (ds DiscoveryService) Route(r *http.Request) (*http.Response, []byte, error) {

	// format URL and service name
	path := strings.TrimPrefix(r.URL.Path, servicePath)
	temp := strings.SplitN(path, "/", 2)

	var serviceName string
	var serviceURL string
	if len(temp) == 1 {
		serviceName = temp[0]
	} else {
		serviceName = temp[0]
		serviceURL = temp[1]
	}

	if _, ok := serviceMap[serviceName]; !ok {
		log.Error("Route Error: invalid service name - " + serviceName)
		return nil, nil, errors.New("Invalid Service Name")
	}
	serviceURL = serviceMap[serviceName] + serviceURL

	// format request body
	body, err := readAllFunc(r.Body)
	if err != nil {
		log.Error("Route Error: " + err.Error())
		return nil, nil, err
	}

	// formate request header
	header := make(map[string]string)
	for key, vals := range r.Header {
		var item string
		for _, val := range vals {
			item += val + ","
		}
		item = item[:len(item)-1]
		header[key] = item
	}
	delete(header, "api-key")
	delete(header, "secret-key")

	// send request
	rsp, body, err := request(serviceURL, r.Method, header, string(body), nil)
	if err != nil {
		log.Error("Route Error: " + err.Error())
	}
	return rsp, body, err
}
