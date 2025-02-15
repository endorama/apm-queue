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

// Package json provides a JSON encoder/decoder.
package json

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

	"github.com/elastic/apm-data/model"
)

func TestJSONMetrics(t *testing.T) {
	rdr := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(rdr))
	meter := mp.Meter("test")

	e, err := meter.Int64Counter("encoded")
	require.NoError(t, err)

	var codec JSON
	RecordEncodedBytes(&codec, e)
	b, err := codec.Encode(model.APMEvent{})
	require.NoError(t, err)

	var rm metricdata.ResourceMetrics
	assert.NoError(t, rdr.Collect(context.Background(), &rm))

	metric := rm.ScopeMetrics[0].Metrics[0]

	metricdatatest.AssertEqual(t, metricdata.Metrics{
		Name: "encoded",
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
			DataPoints: []metricdata.DataPoint[int64]{
				{Value: 2152},
			},
		},
	}, metric, metricdatatest.IgnoreTimestamp())

	d, err := meter.Int64Counter("decoded")
	require.NoError(t, err)

	RecordDecodedBytes(&codec, d)

	require.NoError(t, codec.Decode(b, &model.APMEvent{}))
	assert.NoError(t, rdr.Collect(context.Background(), &rm))

	// Decoded metric
	metric = rm.ScopeMetrics[0].Metrics[1]
	metricdatatest.AssertEqual(t, metricdata.Metrics{
		Name: "decoded",
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
			DataPoints: []metricdata.DataPoint[int64]{
				{Value: 2152},
			},
		},
	}, metric, metricdatatest.IgnoreTimestamp())
}
