package godruid

import (
	"encoding/json"
)

type Aggregation struct {
	Type        string       `json:"type"`
	Name        string       `json:"name,omitempty"`
	FieldName   string       `json:"fieldName,omitempty"`
	FieldNames  []string     `json:"fieldNames,omitempty"`
	FnAggregate string       `json:"fnAggregate,omitempty"`
	FnCombine   string       `json:"fnCombine,omitempty"`
	FnReset     string       `json:"fnReset,omitempty"`
	ByRow       *ByRow       `json:"byRow,omitempty"`
	AggFilter   *Filter      `json:"filter,omitempty"`
	Aggregator  *Aggregation `json:"aggregator,omitempty"`
}

// ---------------------------------
// Options
// ---------------------------------

type AggOption interface {
	apply(*Aggregation)
}

type ByRow bool

func (b ByRow) apply(c *Aggregation) { c.ByRow = &b }

// ---------------------------------
// Constructors
// ---------------------------------

func AggRawJson(rawJson string) Aggregation {
	agg := &Aggregation{}
	json.Unmarshal([]byte(rawJson), agg)
	return *agg
}

func AggCount(name string) Aggregation {
	return Aggregation{
		Type: "count",
		Name: name,
	}
}

func AggLongSum(name, fieldName string) Aggregation {
	return Aggregation{
		Type:      "longSum",
		Name:      name,
		FieldName: fieldName,
	}
}

func AggDoubleSum(name, fieldName string) Aggregation {
	return Aggregation{
		Type:      "doubleSum",
		Name:      name,
		FieldName: fieldName,
	}
}

func AggLongMin(name, fieldName string) Aggregation {
	return Aggregation{
		Type:      "longMin",
		Name:      name,
		FieldName: fieldName,
	}
}

func AggLongMax(name, fieldName string) Aggregation {
	return Aggregation{
		Type:      "longMax",
		Name:      name,
		FieldName: fieldName,
	}
}

func AggDoubleMin(name, fieldName string) Aggregation {
	return Aggregation{
		Type:      "doubleMin",
		Name:      name,
		FieldName: fieldName,
	}
}

func AggDoubleMax(name, fieldName string) Aggregation {
	return Aggregation{
		Type:      "doubleMax",
		Name:      name,
		FieldName: fieldName,
	}
}

func AggJavaScript(name, fnAggregate, fnCombine, fnReset string, fieldNames []string) Aggregation {
	return Aggregation{
		Type:        "javascript",
		Name:        name,
		FieldNames:  fieldNames,
		FnAggregate: fnAggregate,
		FnCombine:   fnCombine,
		FnReset:     fnReset,
	}
}

func AggCardinality(name string, fieldNames []string, options ...AggOption) Aggregation {
	agg := Aggregation{
		Type:       "cardinality",
		Name:       name,
		FieldNames: fieldNames,
	}
	for _, opt := range options {
		opt.apply(&agg)
	}
	return agg
}

func AggHyperUnique(name, fieldName string) Aggregation {
	return Aggregation{
		Type:      "hyperUnique",
		Name:      name,
		FieldName: fieldName,
	}
}

func AggFiltered(aggFilter Filter, aggregation Aggregation) Aggregation {
	return Aggregation{
		Type:       "filtered",
		AggFilter:  &aggFilter,
		Aggregator: &aggregation,
	}
}

func AggThetaSketch(name, fieldName string) Aggregation {
	return Aggregation{
		Type:      "thetaSketch",
		Name:      name,
		FieldName: fieldName,
	}
}
