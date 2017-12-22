package handler

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//CommonHandler shared handler
type CommonHandler struct {
	AllowedMethods []string
}

//ApplyMiddleware apply middleware
func (ch CommonHandler) ApplyMiddleware(next http.Handler) http.Handler {
	return ch.bodyCloser(ch.checkKey(ch.checkMethods(next)))
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
		log.Println(r.Method, "is not allowed")
		return
	})
}

func (ch CommonHandler) bodyCloser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		defer log.Println("Closed Body")
		next.ServeHTTP(w, r)
	})
}

func (ch CommonHandler) checkKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("api-key")
		if apiKey != viper.Get("key.api") {
			w.WriteHeader(http.StatusForbidden)
			log.Println("Forbidden access")
			getIP(r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
