package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setUpMiddleWare() {
	apiKey = "test"
}

func TestApplyMiddleWare(t *testing.T) {
	setUpMiddleWare()
	ch := CommonHandler{}
	ch.AllowedMethods = []string{"GET"}

	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("api-key", "fail")

	rr := httptest.NewRecorder()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	handler := ch.ApplyMiddleware(next)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}

func TestMiddleWareMethodFail(t *testing.T) {
	setUpMiddleWare()
	ch := CommonHandler{}
	ch.AllowedMethods = []string{"GET"}

	req, err := http.NewRequest("POST", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("api-key", "test")

	rr := httptest.NewRecorder()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	handler := ch.ApplyMiddleware(next)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestMiddleWareSuccess(t *testing.T) {
	setUpMiddleWare()
	ch := CommonHandler{}
	ch.AllowedMethods = []string{"GET"}

	req, err := http.NewRequest("GET", "/health-check", strings.NewReader(`{"URL" : ""}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("api-key", "test")

	rr := httptest.NewRecorder()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	handler := ch.ApplyMiddleware(next)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
