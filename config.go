package main

import (
	"encoding/json"
	"os"
)

const configFile = "config.json"

var proxyConfig config

//fields need to be public for (un)marshaling to work
type config struct {
	Port     int    `json:"port" default:"8080"`
	Hostname string `json:"hostname" default:"localhost"`
	Logfile  string `json:"logfile" default:"proxy.log"`
}

func (*config) Load() {
	readConfig, _ := os.ReadFile(configFile)
	json.Unmarshal(readConfig, &proxyConfig)
}

//I added Save for my own education, its not needed.
func (*config) Save() {
	data, _ := json.Marshal(proxyConfig)
	os.WriteFile(configFile, data, os.FileMode(os.O_TRUNC|os.O_CREATE))
}
