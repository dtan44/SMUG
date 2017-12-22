package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var (
	secretKey string
)

func init() {
	secretKey = viper.GetString("key.secret")
}

func validateKey(key string) bool {
	if key == secretKey {
		return true
	}
	return false
}

func readJSONBody(body io.ReadCloser, t interface{}) {
	res, err := ioutil.ReadAll(body)
	if err != nil {
		log.Error("Error readJSONBody: " + err.Error())
	}

	// check if body is empty
	if len(res) == 0 {
		return
	}

	err = json.Unmarshal(res, t)
	if err != nil {
		log.Error("Error readJSONBody: " + err.Error())
	}
}

// https://blog.golang.org/context/userip/userip.go
func getIP(req *http.Request) {

	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		log.Errorf("UserIP: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		log.Errorf("userip: %q is not IP:port", req.RemoteAddr)
		return
	}

	// This will only be defined when site is accessed via non-anonymous proxy
	// and takes precedence over RemoteAddr
	forward := req.Header.Get("X-Forwarded-For")

	log.Infof("IP: %s", ip)
	log.Infof("Port: %s", port)
	log.Infof("Forwarded for: %s", forward)
}
