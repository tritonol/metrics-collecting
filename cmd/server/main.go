package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tritonol/metrics-collecting.git/internal/backup"
	"github.com/tritonol/metrics-collecting.git/internal/routes"
	"github.com/tritonol/metrics-collecting.git/internal/server/config"
	"github.com/tritonol/metrics-collecting.git/internal/storage/memstorage"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()
	storage := memstorage.NewMemStorage()
	
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if cfg.Backup.Restore {
		err := backup.RestoreMetricsFromFile(cfg.Backup.FilePath, storage)
		if err != nil {
			logger.Error("Error restoring metrics:", zap.Error(err))
		}
	}

	go backup.SaveMetricsPeriodically(cfg.Backup.StoreInterval, cfg.Backup.FilePath, storage)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		<-stopChan
		err := backup.SaveMetricsToFile(cfg.Backup.FilePath, storage)
		if err != nil {
			logger.Error("Error save metric:", zap.Error(err))
		}

		logger.Info("Received interrupt signal. Saving data and shutting down")

		os.Exit(0)
	}()

	err := http.ListenAndServe(cfg.Server.Address, routes.MetricRouter(storage, logger))
	if err != nil {
		panic(err)
	}
}