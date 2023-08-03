package main

import (
	"time"

	"github.com/tritonol/metrics-collecting.git/internal/agent/metrics"
	"github.com/tritonol/metrics-collecting.git/internal/agent/request"
)

type MetricRequest interface {
	gatherGauge()
	gatherCounter()
}

const(
	pollInterval = 2 * time.Second
	reportInteravl = 10 * time.Second
	serverAddress = "http://localhost:8080"
)

func main() {
	metrics := metrics.NewMetrics()

	updateTicker := time.NewTicker(pollInterval)
	defer updateTicker.Stop()

	sendTicker := time.NewTicker(reportInteravl)
	defer sendTicker.Stop()

	for {
		select {
		case <-updateTicker.C:
			metrics.CollectCounter()
			metrics.CollectGauge()
		case <-sendTicker.C:
			request.Send(metrics, serverAddress)
		}
	}

}
