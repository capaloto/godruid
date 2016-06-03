package godruid

type Granularity interface{}

type SimpleGran string

type ComplexGran struct {
	Type     string `json:"type"`
	Duration int    `json:"duration,omitempty"`
	Origin   string `json:"origin,omitempty"`
	Period   string `json:"period,omitempty"`
	TimeZone string `json:"timeZone,omitempty"`
}

// ---------------------------------
// Options
// ---------------------------------

type GranOption interface {
	apply(*ComplexGran)
}

type Origin string

func (s Origin) apply(c *ComplexGran) { c.Origin = string(s) }

type TimeZone string

func (s TimeZone) apply(c *ComplexGran) { c.TimeZone = string(s) }

// ---------------------------------
// Contructors
// ---------------------------------

const (
	GranAll        SimpleGran = "all"
	GranNone       SimpleGran = "none"
	GranMinute     SimpleGran = "minute"
	GranFifteenMin SimpleGran = "fifteen_minute"
	GranThirtyMin  SimpleGran = "thirty_minute"
	GranHour       SimpleGran = "hour"
	GranDay        SimpleGran = "day"
)

func GranDuration(duration int, options ...GranOption) Granularity {
	gran := ComplexGran{Type: "duration", Duration: duration}
	for _, opt := range options {
		opt.apply(&gran)
	}
	return gran
}

func GranPeriod(period string, options ...GranOption) Granularity {
	gran := ComplexGran{Type: "period", Period: period}
	for _, opt := range options {
		opt.apply(&gran)
	}
	return gran
}
