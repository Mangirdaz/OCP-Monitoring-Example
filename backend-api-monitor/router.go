package main

import (
	"fmt"
	"net/http"

	"github.com/mangirdaz/ocp-demo/config"
)

func NewRouter() {

	http.HandleFunc("/healthz", Health)
	http.HandleFunc("/readiness", Readiness)
	url := fmt.Sprintf("%s:%s", config.Get("EnvAPIMonIP"), config.Get("EnvAPIMonPort"))

	http.ListenAndServe(url, nil)

}
