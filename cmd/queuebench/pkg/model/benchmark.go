package model

import (
	"fmt"
	"io"
	"math"
	"runtime"
	"strings"
)

// https://go.googlesource.com/proposal/+/master/design/14313-benchmark-format.md
func (r *BenchResult) ToGoBenchmark(w io.Writer) {
	// Go infos
	fmt.Fprintln(w, "goos:", runtime.GOOS)
	fmt.Fprintln(w, "goarch:", runtime.GOARCH)
	fmt.Fprintln(w, "compiler:", runtime.Compiler)
	fmt.Fprintln(w, "cpu-count:", runtime.NumCPU())
	fmt.Fprintln(w, "gomaxprocs:", runtime.GOMAXPROCS(0))
	// Benchmark run configs
	fmt.Fprintln(w, "duration:", r.Meta.Config.Duration)
	fmt.Fprintln(w, "timeout:", r.Meta.Config.Timeout)
	fmt.Fprintln(w, "event-size:", r.Meta.Config.EventSize)
	fmt.Fprintln(w, "partitions:", r.Meta.Config.Partitions)
	// Benchmark metadata
	fmt.Fprintln(w, "# run-id:", r.Meta.RunID)
	fmt.Fprintln(w, "# start-time:", r.Meta.StartTime)
	fmt.Fprintln(w, "# end-time:", r.Meta.EndTime)
	fmt.Fprintln(w, "# delay-min:", r.MinConsumptionDelay)
	fmt.Fprintln(w, "# delay-max:", r.MaxConsumptionDelay)
	fmt.Fprintln(w, "# delay-sum:", r.SumConsumptionDelay)
	fmt.Fprintln(w, "# delay-total-count:", r.ConsumptionDelayTotalCount)
	fmt.Fprintln(w, "# delay-values:", strings.Join(map2string[float64](r.ConsumptionDelay.Values), ","))
	fmt.Fprintln(w, "# delay-counts:", strings.Join(map2string(r.ConsumptionDelay.Counts), ","))
	// The benchmark output line in appropriate format: <name> <iterations> <value> <unit> [<value> <unit>...]
	B2MB := 1000000.0
	line := strings.Builder{}
	line.WriteString(fmt.Sprintf("BenchmarkQueueRun 1 "))
	line.WriteString(fmt.Sprintf("%d produced/op ", r.Produced))
	line.WriteString(fmt.Sprintf("%.2f produced-MB/s ", float64(r.ProducedBytes)/B2MB))
	line.WriteString(fmt.Sprintf("%d consumed/op ", r.Consumed))
	line.WriteString(fmt.Sprintf("%.2f consumed-MB/s ", float64(r.ConsumedBytes)/B2MB))
	line.WriteString(fmt.Sprintf("%d leftover/op ", r.Leftover))
	line.WriteString(fmt.Sprintf("%f p50/op ", P(50, r.ConsumptionDelay)))
	line.WriteString(fmt.Sprintf("%f p90/op ", P(90, r.ConsumptionDelay)))
	line.WriteString(fmt.Sprintf("%f p95/op ", P(95, r.ConsumptionDelay)))
	fmt.Fprintln(w, line.String())
}

func P(p float64, r Histogram) float64 {
	if p < 0 || p > 100 {
		return math.NaN()
	}

	total := 0
	for _, v := range r.Counts {
		total += v
	}
	counts := r.Counts
	boundaries := r.Values

	op := 0.0
	for a, b := range counts {
		cp := float64(b * 100 / total)
		if op+cp > p {
			return boundaries[a]
		}
		op += cp
	}
	return boundaries[len(boundaries)-1]
}

func map2string[T float64 | int](v []T) (s []string) {
	var fmts string
	switch any(v).(type) {
	case []float64:
		fmts = "%f"
	case []int:
		fmts = "%d"
	default:
		fmts = "%+T"
	}

	for _, v := range v {
		s = append(s, fmt.Sprintf(fmts, v))
	}
	return s
}
