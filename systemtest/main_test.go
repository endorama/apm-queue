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

// Package systemtest contains the queue slow system tests
package systemtest

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func testMain(m *testing.M) int {
	var skipDestroy bool
	flag.BoolVar(&skipDestroy, "skip-destroy", false, "do not destroy the provisioned infrastructure after the tests finish")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := ProvisionKafka(ctx, newLocalKafkaConfig()); err != nil {
			return fmt.Errorf("failed to provision kafka: %w", err)
		}
		return nil
	})
	g.Go(func() error {
		if err := ProvisionPubSubLite(ctx, newPubSubLiteConfig()); err != nil {
			return fmt.Errorf("failed to provision pubsublite: %w", err)
		}
		return nil
	})
	if !skipDestroy {
		defer Destroy()
	}
	if err := g.Wait(); err != nil {
		logger().Error(err)
		return 1
	}
	logger().Info("running system tests...")
	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}
