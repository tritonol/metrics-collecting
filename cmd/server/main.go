package main

import (
	"fmt"
	"net/http"

	"github.com/tritonol/metrics-collecting.git/internal/backup"
	"github.com/tritonol/metrics-collecting.git/internal/routes"
	"github.com/tritonol/metrics-collecting.git/internal/server/config"
	"github.com/tritonol/metrics-collecting.git/internal/storage/memstorage"
)

func main() {
	cfg := config.MustLoad()
	storage := memstorage.NewMemStorage()

	if cfg.Backup.Restore {
		err := backup.RestoreMetricsFromFile(cfg.Backup.FilePath, storage)
		if err != nil {
			fmt.Println("Error restoring metrics:", err)
		}
	}

	go backup.SaveMetricsPeriodically(cfg.Backup.StoreInterval, cfg.Backup.FilePath, storage)

	err := http.ListenAndServe(cfg.Server.Address, routes.MetricRouter(storage))
	if err != nil {
		panic(err)
	}
}