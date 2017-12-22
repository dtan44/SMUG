package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"../service"
)

// Global variables
var jsonMarshal func(v interface{}) ([]byte, error)

const (
	registerPath   = "/register/"
	deregisterPath = "/deregister/"
)

func init() {
	jsonMarshal = json.Marshal
}

//Result JSON response body
type Result struct {
	Result string `json:"result,omitempty"`
	Reason string `json:"reason,omitempty"`
}

//ServiceHandler struct
type ServiceHandler struct {
	Registration service.RegistrationInterface
	Discovery    service.DiscoveryInterface
}

type registerBody struct {
	URL string `json:"URL"`
}

//HandleRegister register service
func (sh ServiceHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	serviceName := strings.TrimPrefix(r.URL.Path, registerPath)
	secretKey := r.Header.Get("secret-key")

	// get serviceURL from request body
	requestBody := registerBody{}
	readJSONBody(r.Body, &requestBody)
	serviceURL := requestBody.URL

	var res Result

	if !validateKey(secretKey) {
		// http.Error(w, "Incorrect Key", http.StatusInternalServerError)
		res = Result{"failure", "Incorrect Key"}
	} else {
		err := sh.Registration.Register(serviceName, serviceURL)

		if err != nil {
			res = Result{"failure", err.Error()}
		} else {
			res = Result{"success", ""}
		}
	}

	j, err := jsonMarshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

//HandleDeregister deregister service
func (sh ServiceHandler) HandleDeregister(w http.ResponseWriter, r *http.Request) {
	serviceName := strings.TrimPrefix(r.URL.Path, deregisterPath)
	secretKey := r.Header.Get("secret-key")
	var res Result

	if !validateKey(secretKey) {
		// http.Error(w, "Incorrect Key", http.StatusInternalServerError)
		res = Result{"failure", "Incorrect Key"}
	} else {
		err := sh.Registration.Deregister(serviceName)

		if err != nil {
			res = Result{"failure", err.Error()}
		} else {
			res = Result{"success", ""}
		}
	}

	j, err := jsonMarshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

//ServicesList List of Services
type ServicesList struct {
	Services []string `json:"services"`
}

//HandleList deregister service
func (sh ServiceHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	serviceList := ServicesList{sh.Discovery.List()}

	j, err := jsonMarshal(serviceList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
