package handler

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var apiKey string

func init() {
	apiKey = viper.GetString("key.api")
}

//CommonHandler shared handler
type CommonHandler struct {
	AllowedMethods []string
}

//ApplyMiddleware apply middleware
func (ch CommonHandler) ApplyMiddleware(next http.Handler) http.Handler {
	return ch.closeBody(ch.checkKey(ch.checkMethods(next)))
}

func (ch CommonHandler) checkMethods(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var allowed = false

		for _, m := range ch.AllowedMethods {
			if m == r.Method {
				allowed = true
				break
			}
		}

		if allowed {
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Info(r.Method, "is not allowed")
		getIP(r)
		return
	})
}

func (ch CommonHandler) closeBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			defer r.Body.Close()
		}
		defer log.Info("Closed Body")
		next.ServeHTTP(w, r)
	})
}

func (ch CommonHandler) checkKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("api-key")
		if key != apiKey {
			w.WriteHeader(http.StatusForbidden)
			log.Println("Forbidden access")
			getIP(r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
