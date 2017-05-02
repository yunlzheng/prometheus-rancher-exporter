package main

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// Exporter Sets up all the runtime and metrics
type Exporter struct {
	rancherURL      string
	accessKey       string
	secretKey       string
	hideSys         bool
	agentIP         string
	environmentUUID string
	environmentID   string
	mutex           sync.RWMutex
	gaugeVecs       map[string]*prometheus.GaugeVec
}

// NewExporter creates the metrics we wish to monitor
func newExporter(rancherURL string, accessKey string, secretKey string, hideSys bool, environmentUUID string, agentIP string) *Exporter {
	gaugeVecs := addMetrics()
	environmentID := CurrentEnvironmentID(rancherURL, accessKey, secretKey, environmentUUID)
	return &Exporter{
		gaugeVecs:       gaugeVecs,
		rancherURL:      rancherURL,
		accessKey:       accessKey,
		secretKey:       secretKey,
		hideSys:         hideSys,
		environmentUUID: environmentUUID,
		agentIP:         agentIP,
		environmentID:   environmentID,
	}
}

// CurrentEnvironmentID function, Get current environmentId by uuid from meta-data
func CurrentEnvironmentID(rancherURL string, accessKey string, secretKey string, environmentUUID string) string {
	apiVer := getAPIVersion(rancherURL)
	url := setEndpoint(rancherURL, "projects", apiVer, "")
	var data = new(Data)
	err := getJSON(url, accessKey, secretKey, &data)
	if err != nil {
		log.Error("Error getting JSON from URL ", url)
		return ""
	}

	for _, x := range data.Data {
		log.Infof("scanning current rancher environment %s %s", x.ID, x.Name)
		if x.UUID == environmentUUID {
			return x.ID
		}
	}

	return ""
}
