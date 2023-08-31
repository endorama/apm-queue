package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
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

	fmt.Println("-- production")
	peps := math.Round(eventsPerSecond(float64(data.Produced), data.Duration.Production))
	fmt.Println("produced", data.Produced, "events in", data.Duration.Production)
	dp, err := time.ParseDuration(data.Duration.Production)
	panicErr(err)
	fmt.Printf("production duration %s %f\n", dp, dp.Seconds())
	fmt.Println("produced events per second:", peps)
	fmt.Println("produced bytes per second:", eventsPerSecond(float64(data.ProducedBytes), data.Duration.Production))

	fmt.Println("-- consumption")
	ceps := math.Round(eventsPerSecond(float64(data.Consumed), data.Duration.Consumption))
	fmt.Println("consumed", data.Consumed, "events in", data.Duration.Consumption)
	fmt.Println("consumed events per second:", ceps)
	fmt.Println("produced bytes per second:", eventsPerSecond(float64(data.ConsumedBytes), data.Duration.Consumption))
	dc, err := time.ParseDuration(data.Duration.Consumption)
	panicErr(err)
	fmt.Printf("consumption duration %s %f\n", dc, dc.Seconds())

	fmt.Println("-- error rate")
	fmt.Printf("not consumed: %d\n", data.Leftover)
	fmt.Printf("producer error rate: %d%%\n", data.Leftover/data.Produced*100)

	top := strings.Builder{}
	bottom := strings.Builder{}
	percentage := strings.Builder{}
	percentile := strings.Builder{}

	pp := 0.0
	for i, v := range data.ConsumptionDelay.Bounds {
		top.WriteString(fmt.Sprintf("%d", int(v.Bound)))
		bottom.WriteString(fmt.Sprintf("%d", v.Count))

		a := float64(v.Count) / float64(data.Produced) * 100
		percentage.WriteString(fmt.Sprintf("%f", a))
		percentile.WriteString(fmt.Sprintf("%f", pp+a))
		pp = pp + a
		if i < len(data.ConsumptionDelay.Bounds)-1 {
			top.WriteString(",")
			bottom.WriteString(",")
			percentage.WriteString(",")
			percentile.WriteString(",")
		}
	}

	fmt.Println("-- consumption latency")
	fmt.Printf("bounds: %s\n", top.String())
	fmt.Printf("count : %s\n", bottom.String())
	fmt.Printf("%%     : %s\n", percentage.String())
	fmt.Printf("p     : %s\n", percentile.String())

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
