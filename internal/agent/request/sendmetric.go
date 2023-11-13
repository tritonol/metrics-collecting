package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	jsonstructs "github.com/tritonol/metrics-collecting.git/internal/structs/JSON"
)

type MetricRequest interface {
	CollectCounter() map[string]int64
	CollectGauge() map[string]float64
}

func sendJSONMetrics(serverAddress, mtype, mname string, mvalue interface{}) {
	url := fmt.Sprintf("%s/update/", serverAddress)

	var body jsonstructs.Metrics
	var delta int64
	var value float64

	switch mtype {
	case "gauge":
		value = mvalue.(float64)
	case "counter":
		delta = mvalue.(int64)
	default:
		fmt.Println("Invalid metric type")
		return
	}

	body = jsonstructs.Metrics{
		ID: mname,
		MType: mtype,
		Delta: &delta,
		Value: &value,
	}

	data, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		fmt.Println("Error requesting metrics:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Server returned non-200 status code:", resp.Status)
		return
	}
}

func Send(metricRequest MetricRequest, serverAddress string) {
	gaugeMetrics := metricRequest.CollectGauge()
	counterMetrics := metricRequest.CollectCounter()

	for metric, value := range gaugeMetrics {
		sendJSONMetrics(serverAddress, "gauge", metric, value)
	}

	for metric, value := range counterMetrics {
		sendJSONMetrics(serverAddress, "counter", metric, value)
	}
}