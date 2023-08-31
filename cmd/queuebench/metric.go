// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func sum(dps []metricdata.DataPoint[int64]) (val int64) {
	for _, dp := range dps {
		val += dp.Value
	}
	return val
}

func getSumInt64Metric(instrument string, metric string, rm metricdata.ResourceMetrics) int64 {
	metrics := filterMetrics(instrument, rm.ScopeMetrics)
	if len(metrics) == 0 {
		return 0
	}

	for _, m := range metrics {
		if m.Name == metric {
			return sum(m.Data.(metricdata.Sum[int64]).DataPoints)
		}
	}

	return 0
}

type HistogramBoundCount struct {
	Boundaries []float64 `json:"boundaries"`
	Counts     []int     `json:"count"`
}

type HistogramResult struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Avg   float64 `json:"avg"`
	Sum   float64 `json:"sum"`
	Count uint64  `json:"count"`

	Bounds HistogramBoundCount `json:"bounds"`
}

func getHistogramFloat64Metric(instrument, metric string, rm metricdata.ResourceMetrics) HistogramResult {
	metrics := filterMetrics(instrument, rm.ScopeMetrics)
	if len(metrics) == 0 {
		return HistogramResult{}
	}

	var values metricdata.Histogram[float64]
	for _, m := range metrics {
		if m.Name == metric {
			values = m.Data.(metricdata.Histogram[float64])
		}
	}

	data := values.DataPoints[0]

	getValueOrEmpty := func(value metricdata.Extrema[float64]) float64 {
		if v, ok := value.Value(); ok {
			return v
		}
		return 0.0
	}
	avg := func(a metricdata.Extrema[float64], b metricdata.Extrema[float64]) float64 {
		return (getValueOrEmpty(a) + getValueOrEmpty(b)) / 2
	}
	getBoundsCounts := func(dp metricdata.HistogramDataPoint[float64]) HistogramBoundCount {
		bcs := HistogramBoundCount{}
		for i, b := range dp.Bounds {
			bcs.Boundaries = append(bcs.Boundaries, b)
			bcs.Counts = append(bcs.Counts, int(dp.BucketCounts[i]))
		}
		return bcs
	}

	return HistogramResult{
		Min:    getValueOrEmpty(data.Min),
		Max:    getValueOrEmpty(data.Max),
		Avg:    avg(data.Min, data.Max),
		Sum:    data.Sum,
		Count:  data.Count,
		Bounds: getBoundsCounts(data),
	}
}
