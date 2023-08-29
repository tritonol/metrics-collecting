package backup

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	jsonstructs "github.com/tritonol/metrics-collecting.git/internal/structs/JSON"
)

type metricGetter interface {
	GetAllDataStructed() map[string]jsonstructs.Metrics
	SaveAllDataStructured(metrics map[string]jsonstructs.Metrics) error
}

func SaveMetricsToFile(filePath string, mg metricGetter) error {
	file, err := os.Create(filePath)
	fmt.Println(file.Name())
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(mg.GetAllDataStructed())
}

func SaveMetricsPeriodically(interval int64, filePath string, mg metricGetter) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		err := SaveMetricsToFile(filePath, mg)
		if err != nil {
			fmt.Println("Error saving metrics:", err)
		}
	}
}

func RestoreMetricsFromFile(filePath string, mg metricGetter) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var data map[string]jsonstructs.Metrics

	decoder := json.NewDecoder(file)
	decoder.Decode(&data)

	if err := mg.SaveAllDataStructured(data); err != nil {
		return err
	}

	return nil
}
