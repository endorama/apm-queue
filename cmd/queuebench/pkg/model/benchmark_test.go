package model

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestP(t *testing.T) {
	data := fixture(t, "testdata/1.json")

	assert.Equal(t, 750.0, P(50, data.ConsumptionDelay))
	assert.Equal(t, 1000.0, P(90, data.ConsumptionDelay))
	assert.Equal(t, 2500.0, P(95, data.ConsumptionDelay))
}

func TestOutputBench(t *testing.T) {
	data := fixture(t, "testdata/1.json")

	expected := `goos: linux
goarch: amd64
compiler: gc
cpu-count: 12
gomaxprocs: 12
duration: 300
timeout: 1500
event-size: 1024
partitions: 1
# start-time: 2023-09-01 10:12:45.979967086 +0000 UTC
# end-time: 2023-09-01 10:34:37.712337657 +0000 UTC
# delay-values: 0.000000,5.000000,10.000000,25.000000,50.000000,75.000000,100.000000,250.000000,500.000000,750.000000,1000.000000,2500.000000,5000.000000,7500.000000,10000.000000
# delay-counts: 0,52924,116536,336120,727340,571531,725743,4040638,6899270,6895937,6783832,380527,0,0,0
BenchmarkQueueRunrun-1693563165 1 27530398 produced/op 28464.71 produced-MB/s 27530398 consumed/op 28464.71 consumed-MB/s 0 leftover/op 750.000000 p50/op 1000.000000 p90/op 2500.000000 p95/op 
`

	b := &bytes.Buffer{}
	data.ToGoBenchmark(b)
	require.Equal(t, expected, b.String())
}

func fixture(t *testing.T, name string) BenchResult {
	t.Helper()
	data, err := LoadFromJSON(name)
	require.NoError(t, err)
	return data
}
