package backup

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	jsonstructs "github.com/tritonol/metrics-collecting.git/internal/structs/JSON"
	"go.uber.org/zap"
)

type MetricStorage interface {
	GetAllDataStructed() map[string]jsonstructs.Metrics
	SaveAllDataStructured(metrics map[string]jsonstructs.Metrics) error
}

type BackupManager struct {
	storage      MetricStorage
	filePath     string
	saveInterval time.Duration
	zapLogger    *zap.Logger
}

func NewBackupManager(storage MetricStorage, filePath string, saveInterval time.Duration, zapLogger *zap.Logger) *BackupManager {
	return &BackupManager{
		storage:      storage,
		filePath:     filePath,
		saveInterval: saveInterval,
		zapLogger:    zapLogger,
	}
}

func (bm *BackupManager) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			if err := bm.saveMetricsToFile(); err != nil {
				bm.zapLogger.Error("Error saving metrics before shutdown: ", zap.Error(err))
			} else {
				bm.zapLogger.Info("Metrics was saving before shutdown")
			}
			return
		case <-time.After(bm.saveInterval):
			if err := bm.saveMetricsToFile(); err != nil {
				bm.zapLogger.Error("Error saving metrics: ", zap.Error(err))
			} else {
				bm.zapLogger.Info("Metrics was saving")
			}
		}
	}
}

func (bm *BackupManager) Restore() error {
	file, err := os.Open(bm.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var metrics map[string]jsonstructs.Metrics

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&metrics); err != nil {
		return err
	}

	if err := bm.storage.SaveAllDataStructured(metrics); err != nil {
		return err
	}

	return nil
}

func (bm *BackupManager) saveMetricsToFile() error {
	file, err := os.Create(bm.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	metrics := bm.storage.GetAllDataStructed()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(metrics); err != nil {
		return err
	}

	return nil
}
