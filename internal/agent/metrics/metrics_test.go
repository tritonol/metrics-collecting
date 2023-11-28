package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectGauge(t *testing.T){
	metric := NewMetrics()
	metric.CollectGauge()
	assert.NotEmpty(t, metric.GetGauge())
}

func TestCollectCounter(t *testing.T){
	metric := NewMetrics()
	metric.CollectCounter()
	metrics := metric.GetGauge()
	assert.NotEmpty(t, metrics)

	assert.Equal(t, int64(1), metrics["PollCount"])
}