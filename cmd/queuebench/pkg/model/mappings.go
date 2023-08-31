package model

import "github.com/elastic/go-elasticsearch/v8/typedapi/types"

func ESMapping() *types.TypeMapping {
	return &types.TypeMapping{
		Properties: map[string]types.Property{
			"meta": &types.ObjectProperty{
				Properties: map[string]types.Property{
					"run":        types.NewConstantKeywordProperty(),
					"start_time": types.NewDateProperty(),
					"end_time":   types.NewDateProperty(),
					"config": types.ObjectProperty{
						Properties: map[string]types.Property{
							"event_size": types.NewIntegerNumberProperty(),
							"partitions": types.NewIntegerNumberProperty(),
							"duration":   types.NewFloatNumberProperty(),
							"timeout":    types.NewFloatNumberProperty(),
						},
					},
				},
			},
			"duration": &types.ObjectProperty{
				Properties: map[string]types.Property{
					"total":       types.NewKeywordProperty(),
					"production":  types.NewKeywordProperty(),
					"consumption": types.NewKeywordProperty(),
				},
			},
			"produced":                      types.NewLongNumberProperty(),
			"produced_bytes":                types.NewLongNumberProperty(),
			"consumed":                      types.NewLongNumberProperty(),
			"consumed_bytes":                types.NewLongNumberProperty(),
			"leftover":                      types.NewIntegerNumberProperty(),
			"consumption_delay":             types.NewHistogramProperty(),
			"min_consumption_delay":         types.NewFloatNumberProperty(),
			"max_consumption_delay":         types.NewFloatNumberProperty(),
			"sum_consumption_delay":         types.NewFloatNumberProperty(),
			"consumption_delay_total_count": types.NewIntegerNumberProperty(),
		},
	}
}
