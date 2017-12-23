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
	var res Result

	// get serviceURL from request body
	requestBody := registerBody{}
	err := readJSONBody(r.Body, &requestBody)
	if err != nil {
		res = Result{"failure", err.Error()}
	} else {
		serviceURL := requestBody.URL

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

//HandleList list services
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

//HandleRoute route services
func (sh ServiceHandler) HandleRoute(w http.ResponseWriter, r *http.Request) {
	res, body, err := sh.Discovery.Route(r)
	if err != nil {
		resp := Result{"failure", err.Error()}
		j, err := jsonMarshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(j)
		return
	}

	// add all headers (including multi-valued headers)
	for key, vals := range res.Header {
		val := ""
		for _, item := range vals {
			val += item + ","
		}
		// remove last comma
		val = val[:len(val)-1]

		w.Header().Set(key, val)
	}

	w.WriteHeader(res.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
