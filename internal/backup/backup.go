package backup

import (
	"encoding/json"
	"os"

	jsonstructs "github.com/tritonol/metrics-collecting.git/internal/structs/JSON"
)

type metricGetter interface {
	GetAllDataStructed() map[string]jsonstructs.Metrics
}

func SaveMetricsToFile(filePath string, mg metricGetter, sync bool) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	if sync {
		file.Sync()
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(mg.GetAllDataStructed())
}
