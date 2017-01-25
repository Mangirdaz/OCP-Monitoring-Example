package main

import "github.com/mangirdaz/ocp-mon-demo/config"

type KeyValueStorageConfig struct {
	Ip   string
	Port string
}

type ErrorMessage struct {
	Code    int
	Message string
}

func InitKeyValueStorageConfig() KeyValueStorageConfig {

	var configuration KeyValueStorageConfig
	configuration.Ip = config.Get("EnvKVStorageIp")
	configuration.Port = config.Get("EnvKVStoragePort")

	return configuration
}
