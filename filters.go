package godruid

type Filter struct {
	Type         string           `json:"type"`
	Dimension    string           `json:"dimension,omitempty"`
	Value        interface{}      `json:"value,omitempty"`
	Pattern      string           `json:"pattern,omitempty"`
	Function     string           `json:"function,omitempty"`
	Field        *Filter          `json:"field,omitempty"`
	Fields       []*Filter        `json:"fields,omitempty"`
	ExtractionFn *DimExtractionFn `json:"extractionFn,omitempty"`
	Query        *QuerySearch     `json:"query,omitempty"`
	Values       []string         `json:"values,omitempty"`
	Lower        string           `json:"lower,omitempty"`
	Upper        string           `json:"upper,omitempty"`
	AlphaNumeric *AlphaNumeric    `json:"alphaNumeric,omitempty"`
	LowerStrict  *LowerStrict     `json:"lowerStrict,omitempty"`
	UpperStrict  *UpperStrict     `json:"upperStrict,omitempty"`
}

// ---------------------------------
// Options
// ---------------------------------

type FilterOption interface {
	apply(*Filter)
}

type AlphaNumeric bool

func (b AlphaNumeric) apply(c *Filter) { c.AlphaNumeric = &b }

type Lower string

func (s Lower) apply(c *Filter) { c.Lower = string(s) }

type LowerStrict bool

func (b LowerStrict) apply(c *Filter) { c.LowerStrict = &b }

type Upper string

func (s Upper) apply(c *Filter) { c.Upper = string(s) }

type UpperStrict bool

func (b UpperStrict) apply(c *Filter) { c.UpperStrict = &b }

// ---------------------------------
// Constructors
// ---------------------------------

func FilterSelector(dimension string, value interface{}) *Filter {
	return &Filter{
		Type:      "selector",
		Dimension: dimension,
		Value:     value,
	}
}

func FilterRegex(dimension, pattern string) *Filter {
	return &Filter{
		Type:      "regex",
		Dimension: dimension,
		Pattern:   pattern,
	}
}

func FilterJavaScript(dimension, function string) *Filter {
	return &Filter{
		Type:      "javascript",
		Dimension: dimension,
		Function:  function,
	}
}

func FilterAnd(filters ...*Filter) *Filter {
	return joinFilters(filters, "and")
}

func FilterOr(filters ...*Filter) *Filter {
	return joinFilters(filters, "or")
}

func FilterNot(filter *Filter) *Filter {
	return &Filter{
		Type:  "not",
		Field: filter,
	}
}

func FilterExtraction(dimension, value string, extractionFn *DimExtractionFn) *Filter {
	return &Filter{
		Type:         "extraction",
		Dimension:    dimension,
		Value:        value,
		ExtractionFn: extractionFn,
	}
}

func FilterSearch(dimension string, query *QuerySearch) *Filter {
	return &Filter{
		Type:      "search",
		Dimension: dimension,
		Query:     query,
	}
}

func FilterIn(dimension string, values ...string) *Filter {
	return &Filter{
		Type:      "in",
		Dimension: dimension,
		Values:    values,
	}
}

func FilterBound(dimension string, options ...FilterOption) *Filter {
	filt := &Filter{
		Type:      "bound",
		Dimension: dimension,
	}
	for _, opt := range options {
		opt.apply(filt)
	}
	return filt
}

// ---------------------------------
// Helpers
// ---------------------------------

func joinFilters(filters []*Filter, connector string) *Filter {
	// Remove null filters.
	p := 0
	for _, f := range filters {
		if f != nil {
			filters[p] = f
			p++
		}
	}
	filters = filters[0:p]

	fLen := len(filters)
	if fLen == 0 {
		return nil
	}
	if fLen == 1 {
		return filters[0]
	}

	return &Filter{
		Type:   connector,
		Fields: filters,
	}
}
