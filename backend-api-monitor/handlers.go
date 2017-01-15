package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/mangirdaz/ocp-demo/config"
)

const NotFound = "Not found"

//Index index method for API
func Index(resp http.ResponseWriter, req *http.Request) {

	log.Debug("/ endpoint called")
	resp.Write([]byte("OK"))

}

//Health endpoit for app server
func Readiness(w http.ResponseWriter, r *http.Request) {
	log.Debug("/readiness endpoint called")

	ok := true
	errMsg := ""

	// Check api
	backendURL := fmt.Sprintf("%s:%s/healthz", config.Get("EnvAPIIP"), config.Get("EnvAPIPort"))
	log.Debugf("Checking backend readiness [%s]", backendURL)
	_, err := http.Get(backendURL)

	if err != nil {
		ok = false
		errMsg += fmt.Sprintf("Backend is not ready")
		log.Debugf("Health probe error: %s", err.Error())
	}

	if ok {
		w.Write([]byte("OK"))
	} else {
		// Send 503
		http.Error(w, errMsg, http.StatusServiceUnavailable)
	}
}

//Health endpoit for app server
func Health(w http.ResponseWriter, r *http.Request) {
	log.Debug("/health endpoint called")
	w.Write([]byte("OK"))
}

func response(obj interface{}, prettyPrint bool, err error, resp http.ResponseWriter, req *http.Request) {
	// Check for an error
HAS_ERR:
	if err != nil {

		if err.Error() == NotFound {
			resp.WriteHeader(http.StatusNotFound)
			return
		}

		log.WithFields(log.Fields{
			"error":  err,
			"method": req.Method,
			"url":    req.URL,
		}).Error("request error")

		code := 500
		errMsg := err.Error()
		if strings.Contains(errMsg, "Permission denied") || strings.Contains(errMsg, "ACL not found") {
			code = 403
		}
		resp.WriteHeader(code)
		resp.Write([]byte(err.Error()))
		return
	}

	// Write out the JSON object
	if obj != nil {
		buf, err := marshall(obj, true)
		if err != nil {
			goto HAS_ERR
		}
		resp.Header().Set("Content-Type", "application/json")

		// encoding/json library has a specific bug(feature) to turn empty slices into json null object,
		// let's make an empty array instead
		if string(buf) == "null" {
			buf = []byte("[]")
		}
		resp.Write(buf)
	}
}

// marshall returns a json byte slice, leaving existing json untouched.
func marshall(obj interface{}, pretty bool) ([]byte, error) {

	var js interface{}
	var buf []byte

	// Only check objects that byte slices and strings for valid json
	switch v := obj.(type) {
	case []byte:
		buf = []byte(v)
	case string:
		buf = []byte(v)
	}

	// If we were given a valid json object, return it as-is
	if buf != nil && json.Unmarshal(buf, &js) == nil {
		return buf, nil
	}

	// Otherwise marshall the object into json
	if pretty {
		return json.MarshalIndent(obj, "", "    ")
	}
	return json.Marshal(obj)
}
