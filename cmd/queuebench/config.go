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
	"flag"
	"log"
	"os"
	"time"
)

type config struct {
	broker     string
	duration   time.Duration
	eventSize  int
	output     *os.File
	partitions int
	timeout    time.Duration
	verbose    bool
}

func (c *config) Parse() {
	b := flag.String("broker", "", "Broker bootstrap URL (host:port) to connect to for this benchmark run")
	d := flag.Int("duration", 0, "Duration is seconds of the production phase of the benchmark")
	o := flag.String("output", "", "io.Writer where to write machine readable output for this benchmark run")
	t := flag.String("timeout", "1m", "Timeout for consuming all records. Benchmark will stop regardless of completion.")
	v := flag.Bool("verbose", false, "Enable additional logging")

	flag.Parse()

	if *b == "" {
		log.Fatal("-broker must be set")
	}
	if *d == 0 {
		log.Fatal("-duration must be set and greater than 0")
	}

	output := os.Stdout
	if *o != "" {
		f, err := os.OpenFile(*o, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("-output: cannot open specified file")
		}
		output = f
	}

	timeout, err := time.ParseDuration(*t)
	if err != nil {
		log.Fatalf("cannot parse -timeout '%s' as duration: %s", *t, err)
	}

	c.broker = *b
	c.duration = time.Duration(*d) * time.Second
	c.eventSize = 1024
	c.output = output
	c.partitions = 1
	c.timeout = timeout
	c.verbose = *v
}
