package service

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	healthCheckPath = "/healthcheck"
	apiKey          = "SimpleMicroService"
)

//TODO: Replace temporary cache with database
var (
	serviceMap  map[string]string
	healthCheck func(URL string) bool
	mapPrint    func(map[string]string) string
	request     func(url, httpMethod string,
		headers map[string]string, body string,
		client clientInterface) (*http.Response, []byte, error)
)

func init() {
	serviceMap = make(map[string]string)
	healthCheck = healthCheckURL
	mapPrint = mapToString
	request = sendRequest
}

//RegistrationInterface defines service methods
type RegistrationInterface interface {
	Register(serviceName, serviceURL string) error
	Deregister(serviceName string) error
}

//Request generic interface
type Request interface {
	sendRequest(url, httpMethod string, headers map[string]string, body string, client clientInterface) ([]byte, error)
}

//RegistrationService defines registration service struct
type RegistrationService struct {
}

//Register perform register service
//  serviceName must be unique
//  serviceURL must pass health check
func (rs RegistrationService) Register(serviceName, serviceURL string) error {
	if !healthCheck(serviceURL) {
		return errors.New("URL Health Check Failed")
	}

	if _, ok := serviceMap[serviceName]; !ok {
		serviceMap[serviceName] = serviceURL
		log.Info(mapPrint(serviceMap))
		return nil
	}

	return errors.New("Service Name already Exist")
}

//Deregister perform deregister service
//  serviceName must exist
func (rs RegistrationService) Deregister(serviceName string) error {

	if _, ok := serviceMap[serviceName]; ok {
		delete(serviceMap, serviceName)
		return nil
	}

	return errors.New("Service Name does not Exist")
}

func mapToString(data map[string]string) string {
	jsonString, err := json.Marshal(data)
	if err != nil {
		log.Error("mapToString Error: " + err.Error())
		return ""
	}
	return string(jsonString)
}

//URL must return 200 to GET baseURL/healthcheck
func healthCheckURL(URL string) bool {
	if URL == "" {
		return false
	}

	// initialize request
	// healthcheckURL := URL + healthCheckPath
	// headers := make(map[string]string)
	// headers["api-key"] = apiKey

	// res, _, err := request(healthcheckURL, http.MethodGet, headers, "", nil)
	// if err != nil {
	// 	return false
	// } else if res.StatusCode != http.StatusOK {
	// 	log.Error("healthCheckURL error: Status not OK")
	// 	return false
	// }
	return true
}
