package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/elastic/apm-queue/cmd/queuebench/pkg/model"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func machineoutput(w io.Writer, run int64, realduration, productionduration, consumptionduration time.Duration, start, end time.Time, cfg config, rm metricdata.ResourceMetrics) error {
	totalbytesproduced := getSumInt64Metric("github.com/twmb/franz-go/plugin/kotel", "messaging.kafka.produce_bytes.count", rm)
	totalbytesfetched := getSumInt64Metric("github.com/twmb/franz-go/plugin/kotel", "messaging.kafka.fetch_bytes.count", rm)
	totalproduced := getSumInt64Metric("github.com/elastic/apm-queue/kafka", "producer.messages.produced", rm)
	totalconsumed := getSumInt64Metric("github.com/elastic/apm-queue/kafka", "consumer.messages.fetched", rm)
	delay := getHistogramFloat64Metric("github.com/elastic/apm-queue/kafka", "consumer.messages.delay", rm)

	data := model.BenchResult{
		Meta: model.BenchMeta{
			RunID:     fmt.Sprintf("run-%d", run),
			StartTime: start,
			EndTime:   end,
			Config: model.BenchConfig{
				Duration:   cfg.duration.Seconds(),
				EventSize:  cfg.eventSize,
				Partitions: cfg.partitions,
				Timeout:    cfg.timeout.Seconds(),
			},
		},
		Duration: model.BenchDuration{
			Total:       realduration.Seconds(),
			Production:  productionduration.Seconds(),
			Consumption: consumptionduration.Seconds(),
		},
		Produced:      totalproduced,
		ProducedBytes: totalbytesproduced,
		Consumed:      totalconsumed,
		ConsumedBytes: totalbytesfetched,
		Leftover:      totalproduced - totalconsumed,
		ConsumptionDelay: model.Histogram{
			Values: delay.Bounds.Boundaries,
			Counts: delay.Bounds.Counts,
		},
		MinConsumptionDelay:        delay.Min,
		MaxConsumptionDelay:        delay.Max,
		SumConsumptionDelay:        delay.Sum,
		ConsumptionDelayTotalCount: delay.Count,
	}

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
