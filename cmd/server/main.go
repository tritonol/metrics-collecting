package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	logger.Info("Server strat")

	if cfg.Backup.Restore {
		err := backup.RestoreMetricsFromFile(cfg.Backup.FilePath, storage)
		if err != nil {
			logger.Error("Error restoring metrics:", zap.Error(err))
		}
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	if cfg.Backup.StoreInterval > 0 {
		ticker := time.NewTicker(time.Duration(cfg.Backup.StoreInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := backup.SaveMetricsToFile(cfg.Backup.FilePath, storage); err != nil {
                    logger.Error("Failed to save data to file", zap.Error(err))
                } else {
                    logger.Info("Data saved to file.")
                }
			case <-interrupt:
				if err := backup.SaveMetricsToFile(cfg.Backup.FilePath, storage); err != nil {
                    logger.Error("Failed to save data to file", zap.Error(err))
				} else {
                    logger.Info("Data saved to file and shutdown")
                }

				os.Exit(0)
			}
		}
	}

	err := http.ListenAndServe(cfg.Server.Address, routes.MetricRouter(storage, logger))
	if err != nil {
		panic(err)
	}
}