package config

import (
	"os"

	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/google/uuid"
)

//default values for application
const (
	//api const
	DefaultAPIPort       = "8000"
	DefaultAPIIP         = "0.0.0.0"
	DefaultAPIMonPort    = "8001"
	DefaultAPIMonIP      = "0.0.0.0"
	DefaultAPIServiceURL = "http://api-svc"
	//basic auth enabled or not
	DefaultBasicAuthentication = false

	//libkv
	DefaultStorageBackend = StorageBoltDB
	StorageBoltDB         = "boltdb"
	StorageConsul         = "consul"
	StorageETCD           = "etcd"

	//frontend defaults
	DefaultFEPort = "8001"
	DefaultFEIP   = "0.0.0.0"

	//storage defaults
	DefaultKVStorageIp   = "0.0.0.0"
	DefaultKVStoragePort = "8500"

	//ENV
	EnvAPIPort          = "API_PORT"
	EnvAPIIP            = "API_IP"
	EnvAPIMonPort       = "API_MON_PORT"
	EnvAPIMonIP         = "API_MON_IP"
	EnvBasicAuth        = "API_BASIC_AUTH"
	EnvAPIServiceURL    = "API_SVC"
	EnvDefaultKVBackend = "STORAGE_BACKEND"
	EnvDatabasePath     = "BOLTDB_LOCATION"

	EnvFEPort = "FE_PORT"
	EnvFEIP   = "FE_IP"

	//consul defaults
	EnvKVStorageIp   = "KEYVAL_STORAGE_IP"
	EnvKVStoragePort = "KEYVAL_STORAGE_PORT"
)

//Options structures for application default and configuration
type Options struct {
	Default     string
	Environment string
}

// GenerateID for Note
func GenerateID() (id string) {
	return uuid.New().String()
}

// Get - gets specified variable from either environment or default one
func Get(variable string) string {

	var config = map[string]Options{
		"EnvAPIPort": {
			Default:     DefaultAPIPort,
			Environment: EnvAPIPort,
		},
		"EnvAPIIP": {
			Default:     DefaultAPIIP,
			Environment: EnvAPIIP,
		},
		"EnvBasicAuth": {
			Default:     strconv.FormatBool(DefaultBasicAuthentication),
			Environment: EnvBasicAuth,
		},
		"EnvFEPort": {
			Default:     DefaultFEPort,
			Environment: EnvFEPort,
		},
		"EnvFEIP": {
			Default:     DefaultFEIP,
			Environment: EnvFEIP,
		},
		"EnvKVStorageIp": {
			Default:     DefaultKVStorageIp,
			Environment: EnvKVStorageIp,
		},
		"EnvKVStoragePort": {
			Default:     DefaultKVStoragePort,
			Environment: EnvKVStoragePort,
		},
		"EnvAPIServiceURL": {
			Default:     DefaultAPIServiceURL,
			Environment: EnvAPIServiceURL,
		},
		"EnvDefaultKVBackend": {
			Default:     DefaultStorageBackend,
			Environment: EnvDefaultKVBackend,
		},
		"EnvDatabasePath": {
			Default:     "data/database.db",
			Environment: EnvDatabasePath,
		},
		"EnvAPIMonPort": {
			Default:     DefaultAPIMonPort,
			Environment: EnvAPIMonPort,
		},
		"EnvAPIMonIP": {
			Default:     DefaultAPIMonIP,
			Environment: EnvAPIMonIP,
		},
	}

	for k, v := range config {
		if k == variable {
			if os.Getenv(v.Environment) != "" {
				log.WithFields(log.Fields{
					"key":   k,
					"value": v.Environment,
				}).Debug("config: setting configuration")
				return os.Getenv(v.Environment)
			}
			log.WithFields(log.Fields{
				"key":   k,
				"value": v.Default,
			}).Debug("config: setting configuration")
			return v.Default

		}
	}
	return ""
}
