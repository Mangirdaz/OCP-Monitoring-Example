package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/rakyll/statik/fs"

	"fmt"

	"github.com/mangirdaz/ocp-demo/config"
	_ "github.com/mangirdaz/ocp-demo/frontend/statik"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {

	statikFS, err := fs.New()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("Failed to load statikFS")
	}

	log.Debug(config.Get("EnvAPIIP"))

	url := fmt.Sprintf("%s:%s", config.Get("EnvFEIP"), config.Get("EnvFEPort"))
	log.WithFields(log.Fields{
		"url": url,
	}).Info("frontend: start")

	http.Handle("/", http.FileServer(statikFS))

	//health check
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	//readiness probe
	http.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		ok := true
		errMsg := ""

		// Check api
		backendURL := fmt.Sprintf("%s:%s/healthz", config.Get("EnvAPIServiceURL"), config.Get("EnvAPIPort"))
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
	})

	http.ListenAndServe(url, nil)
}
