package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"./config"
	"./handler"
	"./service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var osSignals chan os.Signal
var serverError chan error

func init() {
	serverError = make(chan error, 1)
	osSignals = make(chan os.Signal, 1)

	signal.Notify(
		osSignals,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGQUIT,
	)
}

func main() {
	var sh handler.ServiceHandler
	sh.Registration = service.RegistrationService{}
	sh.Discovery = service.DiscoveryService{}

	var get handler.CommonHandler
	get.AllowedMethods = []string{http.MethodGet}
	http.Handle("/list", get.ApplyMiddleware(http.HandlerFunc(sh.HandleList)))

	var delete handler.CommonHandler
	delete.AllowedMethods = []string{http.MethodDelete}
	http.Handle("/deregister/", delete.ApplyMiddleware(http.HandlerFunc(sh.HandleDeregister)))

	var put handler.CommonHandler
	put.AllowedMethods = []string{http.MethodPut}
	http.Handle("/register/", put.ApplyMiddleware(http.HandlerFunc(sh.HandleRegister)))

	var route handler.CommonHandler
	route.AllowedMethods = []string{http.MethodGet, http.MethodPost,
		http.MethodPut, http.MethodDelete, http.MethodHead,
		http.MethodConnect, http.MethodOptions, http.MethodPatch,
		http.MethodTrace}
	http.Handle("/service/", route.ApplyMiddleware(http.HandlerFunc(sh.HandleRoute)))

	//TODO: add custom handler for / as a catch all, http has its own default which returns a 404
	//http.Handle("/", mh.BodyCloser(http.HandlerFunc(bye)))

	go runHTTP()

	waitForEvent()

}

func runHTTP() {
	log.Info("Server started")
	log.Info("Listening on Port " + viper.GetString(config.Port))
	serverError <- http.ListenAndServe(":"+viper.GetString(config.Port), nil)
}

func waitForEvent() {

	select {
	case sig := <-osSignals:
		log.Printf("OS shutdown signal:%+v", sig)

	case err := <-serverError:
		log.Println(err)
		panic(err)
	}
	log.Println("Server exit")

}
