package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

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

func machineoutput(w io.Writer, run int64, realduration, productionduration, consumptionduration time.Duration, start, end time.Time, cfg config, rm metricdata.ResourceMetrics) error {
	totalbytesproduced := getSumInt64Metric("github.com/twmb/franz-go/plugin/kotel", "messaging.kafka.produce_bytes.count", rm)
	totalbytesfetched := getSumInt64Metric("github.com/twmb/franz-go/plugin/kotel", "messaging.kafka.fetch_bytes.count", rm)
	totalproduced := getSumInt64Metric("github.com/elastic/apm-queue/kafka", "producer.messages.produced", rm)
	totalconsumed := getSumInt64Metric("github.com/elastic/apm-queue/kafka", "consumer.messages.fetched", rm)
	delay := getHistogramFloat64Metric("github.com/elastic/apm-queue/kafka", "consumer.messages.delay", rm)

	data := ResultData{
		Meta: BenchMeta{
			Run:       fmt.Sprintf("run-%d", run),
			StartTime: start,
			EndTime:   end,
			Config: BenchConfig{
				Duration:   cfg.duration.String(),
				EventSize:  cfg.eventSize,
				Partitions: cfg.partitions,
				Timeout:    cfg.timeout.String(),
			},
		},
		Duration: BenchDuration{
			Total:       realduration.String(),
			Production:  productionduration.String(),
			Consumption: consumptionduration.String(),
		},
		Produced:         totalproduced,
		ProducedBytes:    totalbytesproduced,
		Consumed:         totalconsumed,
		ConsumedBytes:    totalbytesfetched,
		Leftover:         totalproduced - totalconsumed,
		ConsumptionDelay: delay,
	}

	// stats(data)

	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cannot marshal result to json: %w", err)
	}

	log.Println("writing machine readable results")
	if _, err := w.Write(b); err != nil {
		log.Panicf("cannot write machine readable results to io.Writer: %s", err)
	}
	w.Write([]byte("\n"))

	return nil
}

// func stats(data ResultData) {
// 	productionduration, err := time.ParseDuration(data.Duration.Production)
// 	if err != nil {
// 		panic(err)
// 	}

// 	producedBytesPerSecond := data.Produced / int64(productionduration)
// 	log.Println("produced_bytes_per_second:", producedBytesPerSecond)

//   consumedBytesPerSecond := data.Consumed /

// }
