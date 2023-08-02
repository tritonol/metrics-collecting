package request

import (
	"fmt"
	"net/http"
)

type MetricRequest interface {
	CollectCounter() map[string]int64
	CollectGauge() map[string]float64
}

func Send(metricRequest MetricRequest, serverAddress string) {
	gaugeMetrics := metricRequest.CollectGauge()
	counterMetrics := metricRequest.CollectCounter()

	for metric, value := range gaugeMetrics {
		url := fmt.Sprintf("%s/update/gauge/%s/%f", serverAddress, metric, value)
		
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}
		req.Header.Set("Content-Type", "text/plain")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error sending metrics:", err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Server returned non-200 status code:", resp.Status)
		}
	}

	for metric, value := range counterMetrics {
		url := fmt.Sprintf("%s/update/counter/%s/%d", serverAddress, metric, value)
		
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}
		req.Header.Set("Content-Type", "text/plain")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error sending metrics:", err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Server returned non-200 status code:", resp.Status)
		}
	}
}
