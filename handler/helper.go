package handler

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var (
	secretKey     string
	readAllFunc   func(r io.Reader) ([]byte, error)
	jsonUnmarshal func(data []byte, v interface{}) error
)

func init() {
	secretKey = viper.GetString("key.secret")
	readAllFunc = ioutil.ReadAll
	jsonUnmarshal = json.Unmarshal
}

func validateKey(key string) bool {
	if key == secretKey {
		return true
	}
	return false
}

func readJSONBody(body io.ReadCloser, t interface{}) error {
	res, err := readAllFunc(body)
	if err != nil {
		log.Error(" readJSONBody Error: " + err.Error())
		return err
	}

	// check if body is empty
	if len(res) == 0 {
		log.Error("readJSONBody Error: Length 0")
		return errors.New("Length 0")
	}

	err = jsonUnmarshal(res, t)
	if err != nil {
		log.Error("readJSONBody Error: " + err.Error())
		return err
	}

	return nil
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
