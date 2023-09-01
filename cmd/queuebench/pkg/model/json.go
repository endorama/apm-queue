package model

import (
	"encoding/json"
	"fmt"
)

func (r *BenchResult) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot marshal result to json: %w", err)
	}

	return b, nil
}
