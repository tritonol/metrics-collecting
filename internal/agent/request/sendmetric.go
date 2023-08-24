package request

import (
	"fmt"
	"net/http"
)

type MetricRequest interface {
	CollectCounter() map[string]int64
	CollectGauge() map[string]float64
}

func sendMetric(serverAddress, metricType, metricName string, metricValue interface{}) {
	url := fmt.Sprintf("%s/update/%s/%s/%v", serverAddress, metricType, metricName, metricValue)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending metrics:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Server returned non-200 status code:", resp.Status)
	}
}

func Send(metricRequest MetricRequest, serverAddress string) {
	gaugeMetrics := metricRequest.CollectGauge()
	counterMetrics := metricRequest.CollectCounter()

	for metric, value := range gaugeMetrics {
		sendMetric(serverAddress, "gauge", metric, value)
	}

	for metric, value := range counterMetrics {
		sendMetric(serverAddress, "counter", metric, value)
	}
}