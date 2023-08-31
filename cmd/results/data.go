package main

import "time"

type HistogramBoundCount struct {
	Bound float64 `json:"bound"`
	Count uint64  `json:"count"`
}

type HistogramResult struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Avg   float64 `json:"avg"`
	Sum   float64 `json:"sum"`
	Count uint64  `json:"count"`

	Bounds []HistogramBoundCount `json:"bounds"`
}

type BenchConfig struct {
	EventSize  int    `json:"event_size"`
	Partitions int    `json:"partitions"`
	Duration   string `json:"duration"`
	Timeout    string `json:"timeout"`
}

type BenchMeta struct {
	Run       string      `json:"run_id"`
	StartTime time.Time   `json:"start_time"`
	EndTime   time.Time   `json:"end_time"`
	Config    BenchConfig `json:"config"`
}

type BenchDuration struct {
	Total       string `json:"total"`
	Production  string `json:"production"`
	Consumption string `json:"consumption"`
}

type ResultData struct {
	Meta             BenchMeta       `json:"meta"`
	Duration         BenchDuration   `json:"duration"`
	Produced         int64           `json:"produced"`
	ProducedBytes    int64           `json:"produced_bytes"`
	Consumed         int64           `json:"consumed"`
	ConsumedBytes    int64           `json:"consumed_bytes"`
	Leftover         int64           `json:"leftover"`
	ConsumptionDelay HistogramResult `json:"consumption_delay"`
}

 
 
 
 
 
 

 

