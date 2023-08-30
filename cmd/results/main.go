package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"time"
)

func main() {
	fmt.Println(len(os.Args), os.Args)
	file := os.Args[1]

	fp, err := os.Open(file)
	panicErr(err)

	b, err := io.ReadAll(fp)
	panicErr(err)

	var data ResultData
	if err := json.Unmarshal(b, &data); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", data)

	peps := math.Round(eventsPerSecond(float64(data.Produced), data.Duration.Production))
	fmt.Println("produced", data.Produced, "events in", data.Duration.Production)
	fmt.Println("produced events per second:", peps)
	fmt.Println("produced bytes per second:", eventsPerSecond(float64(data.ProducedBytes), data.Duration.Production))
	ceps := math.Round(eventsPerSecond(float64(data.Consumed), data.Duration.Consumption))
	fmt.Println("consumed events per second:", ceps)
	fmt.Println("produced bytes per second:", eventsPerSecond(float64(data.ConsumedBytes), data.Duration.Consumption))

	fmt.Println("producer error rate:")
	fmt.Println("consumption latency p50:")
	fmt.Println("consumption latency p90:")
	fmt.Println("consumption latency p99:")
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func eventsPerSecond(a float64, v string) float64 {
	t, err := time.ParseDuration(v)
	if err != nil {
		panic(err)
	}
	return a / t.Seconds()
}
