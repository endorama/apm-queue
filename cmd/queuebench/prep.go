package main

import (
	"log"

	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/zap"
)

func prepLogger(verbose bool) *zap.Logger {
	var logger *zap.Logger
	var err error
	if verbose {
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatalf("cannot create zap logger: %s", err)
		}
	} else {
		logger = zap.NewNop()
	}
	return logger
}

// taken from https://github.com/elastic/apm-agent-go/blob/8041dd706d18cb72693f15534c54b390050f0a54/module/apmotel/gatherer_config.go#L25
var customHistogramBoundaries = []float64{
	0.00390625, 0.00552427, 0.0078125, 0.0110485, 0.015625, 0.0220971, 0.03125,
	0.0441942, 0.0625, 0.0883883, 0.125, 0.176777, 0.25, 0.353553, 0.5, 0.707107,
	1, 1.41421, 2, 2.82843, 4, 5.65685, 8, 11.3137, 16, 22.6274, 32, 45.2548, 64,
	90.5097, 128, 181.019, 256, 362.039, 512, 724.077, 1024, 1448.15, 2048,
	2896.31, 4096, 5792.62, 8192, 11585.2, 16384, 23170.5, 32768, 46341.0, 65536,
	92681.9, 131072,
}

func prepMeterProvider() *sdkmetric.MeterProvider {
	rdr := sdkmetric.NewManualReader()
	return sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(rdr),
		sdkmetric.WithView(
			sdkmetric.NewView(
				sdkmetric.Instrument{
					Name:  "consumer.messages.delay",
					Scope: instrumentation.Scope{Name: "github.com/elastic/apm-queue/kafka"},
				},
				sdkmetric.Stream{
					Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
						Boundaries: customHistogramBoundaries,
					},
				},
			),
		),
	)
}
