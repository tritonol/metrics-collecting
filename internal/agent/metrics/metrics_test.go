package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectGauge(t *testing.T){
	assert.NotEmpty(t, NewMetrics().CollectGauge())
}

func TestCollectCounter(t *testing.T){
	metrics := NewMetrics().CollectCounter()

	assert.NotEmpty(t, metrics)

	assert.Equal(t, int64(1), metrics["PollCount"])
}