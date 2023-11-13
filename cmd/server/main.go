package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/tritonol/metrics-collecting.git/internal/backup"
	"github.com/tritonol/metrics-collecting.git/internal/routes"
	"github.com/tritonol/metrics-collecting.git/internal/server/config"
	"github.com/tritonol/metrics-collecting.git/internal/storage/memstorage"
	"github.com/tritonol/metrics-collecting.git/internal/storage/pgstorage"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	storage := memstorage.NewMemStorage()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	db, err := pgstorage.NewPg(ctx, cfg.Db.ConnString)
	if err != nil {
		logger.Error("Can`t connect db", zap.Error(err))
	}

	server := &http.Server{Addr: cfg.Server.Address, Handler: routes.MetricRouter(ctx, db, storage, logger)}

	logger.Info("Server strat")

	var wg sync.WaitGroup

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	backupManager := backup.NewBackupManager(storage, cfg.Backup.FilePath, time.Duration(cfg.Backup.StoreInterval)*time.Second, logger)
	if cfg.Backup.Restore {
		if err := backupManager.Restore(); err != nil {
			logger.Error("Error restoring metrics:", zap.Error(err))
		} else {
			logger.Info("Data was restored")
		}
	}

	wg.Add(1)
	go backupManager.Start(ctx, &wg)

	go func() {
		<-sig
		err := server.Shutdown(ctx)
		if err != nil {
			logger.Fatal("", zap.Error(err))
		}
		cancel()
	}()

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal("", zap.Error(err))
		cancel()
	}
	wg.Wait()
}
