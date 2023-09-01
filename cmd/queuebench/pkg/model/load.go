package model

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func LoadFromJSON(name string) (BenchResult, error) {
	fp, err := os.Open(name)
	if err != nil {
		return BenchResult{}, fmt.Errorf("cannot open file: %w", err)
	}
	b, err := io.ReadAll(fp)
	if err != nil {
		return BenchResult{}, fmt.Errorf("cannot read file: %w", err)
	}

	var data BenchResult
	err = json.Unmarshal(b, &data)
	if err != nil {
		return BenchResult{}, fmt.Errorf("cannot unmarshal content: %w", err)
	}

	return data, nil
}
