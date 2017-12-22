package service

import (
	"net/http"
	"testing"
)

func setupServiceRegister() {
	serviceMap = make(map[string]string)
	healthCheck = healthCheckURL
	mapPrint = mapToString
	request = func(url, httpMethod string,
		headers map[string]string, body string,
		client clientInterface) (*http.Response, []byte, error) {
		res := http.Response{}
		res.StatusCode = http.StatusOK

		return &res, []byte{}, nil
	}
}

func TestRegisterURLFail(t *testing.T) {
	setupServiceRegister()
	var rs RegistrationService

	errorText := "URL Health Check Failed"

	// invalid URL
	err := rs.Register("", "")
	if err.Error() != errorText {
		t.Errorf("service returned unexpected error: got %v want %v",
			err.Error(), errorText)
	}
}

func TestRegisterMapFail(t *testing.T) {
	setupServiceRegister()
	var rs RegistrationService

	serviceMap = make(map[string]string)
	serviceMap["test"] = "test"
	errorText := "Service Name already Exist"

	err := rs.Register("test", "test")
	if err.Error() != errorText {
		t.Errorf("service returned unexpected error: got %v want %v",
			err.Error(), errorText)
	}
}

func TestRegisterSuccess(t *testing.T) {
	setupServiceRegister()
	var rs RegistrationService

	// valid URL
	err := rs.Register("test", "test")
	if err != nil {
		t.Errorf("Failed to register service")
	}
}

func TestDeregisterFail(t *testing.T) {
	setupServiceRegister()
	var rs RegistrationService

	errorText := "Service Name does not Exist"

	err := rs.Deregister("test")
	if err.Error() != errorText {
		t.Errorf("service returned unexpected error: got %v want %v",
			err.Error(), errorText)
	}
}

func TestDeregisterSuccess(t *testing.T) {
	setupServiceRegister()
	var rs RegistrationService

	serviceMap["test"] = "test"

	err := rs.Deregister("test")
	if err != nil {
		t.Errorf("Failed to register service")
	}
}
