package main

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type TrafficRouter struct {
	AvailableOrigins []string `envconfig:"AVAILABLE_ORIGINS"`
	HTTPPort         string   `envconfig:"HTTP_PORT"`
	CurrentOrigin    string
}

func (tr *TrafficRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, " ", r.Method, " ", r.URL)

	if strings.Contains(r.URL.Path, "/switch") {
		if tr.CurrentOrigin == tr.AvailableOrigins[0] {
			tr.CurrentOrigin = tr.AvailableOrigins[1]
		} else {
			tr.CurrentOrigin = tr.AvailableOrigins[0]
		}
		return
	} else {
		log.Println("redirecting to ", tr.CurrentOrigin+r.URL.Path)
		resp, err := http.Get(tr.CurrentOrigin + r.URL.Path)
		if err != nil {
			panic(err)
		} else {
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
		}
	}
}

func main() {
	var tr TrafficRouter
	if err := envconfig.Process("", &tr); err != nil {
		log.Fatalln(err)
	}
	http.Handle("/", &tr)
	http.ListenAndServe(":"+tr.HTTPPort, nil)
}
