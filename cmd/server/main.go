package main

import (
	"net/http"
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

	backupManager := backup.NewBackupManager(storage, cfg.Backup.FilePath, time.Duration(cfg.Backup.StoreInterval)*time.Second, logger)
	if cfg.Backup.Restore {
		if err := backupManager.Restore(); err != nil {
			logger.Error("Error restoring metrics:", zap.Error(err))
		}
		logger.Info("Data was restored")
	}

	go backupManager.Start()

	err := http.ListenAndServe(cfg.Server.Address, routes.MetricRouter(storage, logger))
	if err != nil {
		panic(err)
	}
}
