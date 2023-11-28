package main

import (
	"time"

	"github.com/tritonol/metrics-collecting.git/internal/agent/config"
	metr "github.com/tritonol/metrics-collecting.git/internal/agent/metrics"
	"github.com/tritonol/metrics-collecting.git/internal/agent/request"
	workerpool "github.com/tritonol/metrics-collecting.git/internal/agent/worker"
)

func main() {
	cfg := config.MustLoad()

	metrics := metr.NewMetrics()

	updateTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer updateTicker.Stop()

	sendTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	defer sendTicker.Stop()

	workerPool := workerpool.NewWorkerPool(cfg.RateLimit)

	go func() {
		for range updateTicker.C {
			metrics.CollectCounter()
			metrics.CollectGauge()
			metrics.CollectAdditionalGauge()
		}
	}()

	go func() {
		for range sendTicker.C {
			workerPool.Submit(func() {
				request.SendBatch(metrics, cfg.Address, cfg.Key)
			})
		}
	}()

	select {}
}
