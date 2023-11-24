package main

import (
	"time"

	"github.com/tritonol/metrics-collecting.git/internal/agent/config"
	"github.com/tritonol/metrics-collecting.git/internal/agent/metrics"
	"github.com/tritonol/metrics-collecting.git/internal/agent/request"
)

func main() {
	cfg := config.MustLoad()

	metrics := metrics.NewMetrics()

	updateTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer updateTicker.Stop()

	sendTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	defer sendTicker.Stop()

	for {
		select {
		case <-updateTicker.C:
			metrics.CollectCounter()
			metrics.CollectGauge()
		case <-sendTicker.C:
			request.SendBatch(metrics, cfg.Address, cfg.Key)
		}
	}

}
