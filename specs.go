package godruid

// Defines some small spec like structs here.

// ---------------------------------
// LimitSpec
// ---------------------------------

type Limit struct {
	Type    string   `json:"type"`
	Limit   int      `json:"limit"`
	Columns []Column `json:"columns,omitempty"`
}

const (
	LimitAsc  = "Ascending"
	LimitDesc = "Descending"
)

type Column struct {
	Dimension string `json:"dimension"`
	Direction string `json:"direction,omitempty"`
}

func LimitDefault(limit int, columns ...Column) *Limit {
	return &Limit{
		Type:    "default",
		Limit:   limit,
		Columns: columns,
	}
}

// ---------------------------------
// PagingSpec
// ---------------------------------

type PagingSpec struct {
	PagingIdentifiers PagingIdentifiers `json:"pagingIdentifiers"`
	Threshold         int               `json:"threshold"`
}

type PagingIdentifiers interface{}

type PagingIdEmpty struct{}

// ---------------------------------
// SearchQuerySpec
// ---------------------------------

type SearchQuery struct {
	Type          string         `json:"type"`
	Value         string         `json:"value,omitempty"`
	Values        []string       `json:"values,omitempty"`
	CaseSensitive *CaseSensitive `json:"caseSensitive,omitempty"`
}

// Options

type SearchQueryOption interface {
	apply(*SearchQuery)
}

type CaseSensitive bool

func (b CaseSensitive) apply(c *SearchQuery) { c.CaseSensitive = &b }

// Constructors

func SearchQueryInsensitiveContains(value string) *SearchQuery {
	return &SearchQuery{
		Type:  "insensitive_contains",
		Value: value,
	}
}

func SearchQueryFragmentSearch(values []string, options ...SearchQueryOption) *SearchQuery {
	query := SearchQuery{Type: "fragment", Values: values}
	for _, opt := range options {
		opt.apply(&query)
	}
	return &query
}

func SearchQueryContains(value string, options ...SearchQueryOption) *SearchQuery {
	query := SearchQuery{Type: "contains", Value: value}
	for _, opt := range options {
		opt.apply(&query)
	}
	return &query
}

// ---------------------------------
// ToInclude
// ---------------------------------

type ToInclude struct {
	Type    string   `json:"type"`
	Columns []string `json:"columns,omitempty"`
}

var (
	ToIncludeAll  = &ToInclude{Type: "All"}
	ToIncludeNone = &ToInclude{Type: "None"}
)

func ToIncludeList(columns []string) *ToInclude {
	return &ToInclude{
		Type:    "list",
		Columns: columns,
	}
}

// ---------------------------------
// TopNMetricSpec
// ---------------------------------

type TopNMetric struct {
	Type         string      `json:"type"`
	Metric       interface{} `json:"metric,omitempty"`
	PreviousStop string      `json:"previousStop,omitempty"`
}

// options

type TopNMetricOption interface {
	apply(*TopNMetric)
}

type PreviousStop string

func (s PreviousStop) apply(c *TopNMetric) { c.PreviousStop = string(s) }

// constructors

func TopNMetricNumeric(metric string) *TopNMetric {
	return &TopNMetric{
		Type:   "numeric",
		Metric: metric,
	}
}

func TopNMetricLexicographic(options ...TopNMetricOption) *TopNMetric {
	tnm := TopNMetric{
		Type: "lexicographic",
	}
	for _, opt := range options {
		opt.apply(&tnm)
	}
	return &tnm
}

func TopNMetricAlphaNumeric(options ...TopNMetricOption) *TopNMetric {
	tnm := TopNMetric{
		Type: "alphaNumeric",
	}
	for _, opt := range options {
		opt.apply(&tnm)
	}
	return &tnm
}

func TopNMetricInverted(metric *TopNMetric) *TopNMetric {
	return &TopNMetric{
		Type:   "inverted",
		Metric: metric,
	}
}

// ---------------------------------
// SearchSortSpec
// ---------------------------------

type SearchSort struct {
	Type string `json:"type"`
}

var (
	SearchSortLexicographic = &SearchSort{Type: "lexicographic"}
	SearchSortStrlen        = &SearchSort{Type: "strlen"}
)

// ---------------------------------
// Lookup
// ---------------------------------

type Lookup struct {
	Type       string            `json:"type"`
	Map        map[string]string `json:"map,omitempty"`
	Namespace  string            `json:"namespace,omitempty"`
	IsOneToOne *IsOneToOne       `json:"isOneToOne,omitempty"`
}

// Options

type LookupOption interface {
	apply(*Lookup)
}

type IsOneToOne bool

func (b IsOneToOne) apply(c *Lookup) { c.IsOneToOne = &b }

// Constructors

func LookupNamespace(namespace string, options ...LookupOption) *Lookup {
	lu := Lookup{
		Type:      "namespace",
		Namespace: namespace,
	}
	for _, opt := range options {
		opt.apply(&lu)
	}
	return &lu
}

func LookupMap(lookupmap map[string]string, options ...LookupOption) *Lookup {
	lu := Lookup{
		Type: "map",
		Map:  lookupmap,
	}
	for _, opt := range options {
		opt.apply(&lu)
	}
	return &lu
}
