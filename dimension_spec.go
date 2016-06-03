package godruid

type DimSpec interface{}

type Dimension struct {
	Type                    string              `json:"type"`
	Dimension               string              `json:"dimension"`
	OutputName              string              `json:"outputName"`
	DimExtractionFn         *DimExtractionFn    `json:"dimExtractionFn,omitempty"`
	RetainMissingValue      *RetainMissingValue `json:"retainMissingValue,omitempty"`
	ReplaceMissingValueWith string              `json:"replaceMissingValueWith,omitempty"`
	Lookup                  *Lookup             `json:"lookup,omitempty"`
	Name                    string              `json:"name,omitempty"`
	Optimize                *Optimize           `json:"optimize,omitempty"`
}

type DimExtractionFn struct {
	Type                    string               `json:"type"`
	Expr                    string               `json:"expr,omitempty"`
	ReplaceMissingValue     *ReplaceMissingValue `json:"replaceMissingValue,omitempty"`
	ReplaceMissingValueWith string               `json:"replaceMissingValueWith,omitempty"`
	Query                   *SearchQuery         `json:"query,omitempty"`
	Index                   int                  `json:"index,omitempty"`
	Length                  int                  `json:"length,omitempty"`
	Format                  string               `json:"format,omitempty"`
	TimeZone                string               `json:"timeZone,omitempty"`
	Locale                  string               `json:"locale,omitempty"`
	TimeFormat              string               `json:"timeFormat,omitempty"`
	ResultFormat            string               `json:"resultFormat,omitempty"`
	Function                string               `json:"function,omitempty"`
	Injective               *Injective           `json:"injective,omitempty"`
	RetainMissingValue      *RetainMissingValue  `json:"retainMissingValue,omitempty"`
	Lookup                  *Lookup              `json:"lookup,omitempty"`
	ExtractionFns           []*DimExtractionFn   `json:"extractionFns,omitempty"`
	Delegate                *DimSpec             `json:"delegate,omitempty"`
	Values                  []string             `json:"values,omitempty"`
	IsWhiteList             *IsWhiteList         `json:"isWhitelist,omitempty"`
	Pattern                 string               `json:"pattern,omitempty"`
	Optimize                *Optimize            `json:"optimize,omitempty"`
}

// ---------------------------------
// Options
// ---------------------------------

type DimOption interface {
	applyDim(*Dimension)
}

type DimExFnOption interface {
	apply(*DimExtractionFn)
}

type Injective bool

func (b Injective) apply(c *DimExtractionFn) { c.Injective = &b }

type IsWhiteList bool

func (b IsWhiteList) apply(c *DimExtractionFn) { c.IsWhiteList = &b }

type Length int

func (i Length) apply(c *DimExtractionFn) { c.Length = int(i) }

type ReplaceMissingValue bool

func (b ReplaceMissingValue) apply(c *DimExtractionFn) { c.ReplaceMissingValue = &b }

type ReplaceMissingValueWith string

func (s ReplaceMissingValueWith) applyDim(c *Dimension)    { c.ReplaceMissingValueWith = string(s) }
func (s ReplaceMissingValueWith) apply(c *DimExtractionFn) { c.ReplaceMissingValueWith = string(s) }

type RetainMissingValue bool

func (b RetainMissingValue) applyDim(c *Dimension)    { c.RetainMissingValue = &b }
func (b RetainMissingValue) apply(c *DimExtractionFn) { c.RetainMissingValue = &b }

type Optimize bool

func (b Optimize) applyDim(c *Dimension)    { c.Optimize = &b }
func (b Optimize) apply(c *DimExtractionFn) { c.Optimize = &b }

// ---------------------------------
// Dimension Constructors
// ---------------------------------

func DimDefault(dimension, outputName string) DimSpec {
	return &Dimension{
		Type:       "default",
		Dimension:  dimension,
		OutputName: outputName,
	}
}

func DimExtraction(dimension, outputName string, fn *DimExtractionFn) DimSpec {
	return &Dimension{
		Type:            "extraction",
		Dimension:       dimension,
		OutputName:      outputName,
		DimExtractionFn: fn,
	}
}

func DimLookupMap(dimension, outputName string, lookup *Lookup, options ...DimOption) DimSpec {
	dim := &Dimension{
		Type:       "lookup",
		Dimension:  dimension,
		OutputName: outputName,
		Lookup:     lookup,
	}
	for _, opt := range options {
		opt.applyDim(dim)
	}
	return dim
}

func DimLookupNamespace(dimension, outputName, namespace string) DimSpec {
	return &Dimension{
		Type:       "lookup",
		Dimension:  dimension,
		OutputName: outputName,
		Name:       namespace,
	}
}

// ---------------------------------
// Extraction Function Constructors
// ---------------------------------

func DimExFnRegex(expr string, options ...DimExFnOption) *DimExtractionFn {
	exFn := &DimExtractionFn{
		Type: "regex",
		Expr: expr,
	}
	for _, opt := range options {
		opt.apply(exFn)
	}
	return exFn
}

func DimExFnPartial(expr string) *DimExtractionFn {
	return &DimExtractionFn{
		Type: "partial",
		Expr: expr,
	}
}

func DimExFnSearchQuerySpec(query *SearchQuery) *DimExtractionFn {
	return &DimExtractionFn{
		Type:  "searchQuery",
		Query: query,
	}
}

func DimExFnSubstringQuerySpec(index int, options ...DimExFnOption) *DimExtractionFn {
	exFn := &DimExtractionFn{
		Type:  "searchQuery",
		Index: index,
	}
	for _, opt := range options {
		opt.apply(exFn)
	}
	return exFn
}

func DimExFnTimeFormat(format, timeZone, locale string) *DimExtractionFn {
	return &DimExtractionFn{
		Type:     "timeFormat",
		Format:   format,
		TimeZone: timeZone,
		Locale:   locale,
	}
}

func DimExFnTime(timeFormat, resultFormat string) *DimExtractionFn {
	return &DimExtractionFn{
		Type:         "time",
		TimeFormat:   timeFormat,
		ResultFormat: resultFormat,
	}
}

func DimExFnJavascript(function string, options ...DimExFnOption) *DimExtractionFn {
	exFn := &DimExtractionFn{
		Type:     "javascript",
		Function: function,
	}
	for _, opt := range options {
		opt.apply(exFn)
	}
	return exFn
}

func DimExFnLookup(lookup *Lookup, options ...DimExFnOption) *DimExtractionFn {
	exFn := &DimExtractionFn{
		Type:   "lookup",
		Lookup: lookup,
	}
	for _, opt := range options {
		opt.apply(exFn)
	}
	return exFn
}

func DimExFnCascade(extractionFns ...*DimExtractionFn) *DimExtractionFn {
	return &DimExtractionFn{
		Type:          "cascade",
		ExtractionFns: extractionFns,
	}
}

func DimExFnStringFormat(format string) *DimExtractionFn {
	return &DimExtractionFn{
		Type:   "stringFormat",
		Format: format,
	}
}

func DimExFnListFiltered(delegate *DimSpec, values []string, options ...DimExFnOption) *DimExtractionFn {
	exFn := &DimExtractionFn{
		Type:     "listFiltered",
		Delegate: delegate,
		Values:   values,
	}
	for _, opt := range options {
		opt.apply(exFn)
	}
	return exFn
}

func DimExFnRegexFiltered(delegate *DimSpec, pattern string) *DimExtractionFn {
	return &DimExtractionFn{
		Type:     "listFiltered",
		Delegate: delegate,
		Pattern:  pattern,
	}
}

func DimExFnUpper(options ...DimExFnOption) *DimExtractionFn {
	exFn := &DimExtractionFn{
		Type: "upper",
	}
	for _, opt := range options {
		opt.apply(exFn)
	}
	return exFn
}

func DimExFnLower(options ...DimExFnOption) *DimExtractionFn {
	exFn := &DimExtractionFn{
		Type: "lower",
	}
	for _, opt := range options {
		opt.apply(exFn)
	}
	return exFn
}
