package godruid

import (
	"fmt"
	"syscall"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGroupby(t *testing.T) {
	Convey("TestGroupby", t, func() {
		broker, found := syscall.Getenv("DruidBroker")
		appId, found2 := syscall.Getenv("TestAppId")
		So(found && found2, ShouldEqual, true)

		query := &QueryGroupBy{
			DataSource:   "events_agg",
			Intervals:    []string{"2016-05-01T00:00/2016-05-01T01"},
			Granularity:  GranAll,
			Filter:       FilterAnd(FilterSelector("app_id", appId), FilterSelector("attribution_network_key", nil), nil),
			LimitSpec:    LimitDefault(5, Column{Dimension: "revenue", Direction: LimitDesc}),
			Dimensions:   []DimSpec{"attribution_network"},
			Aggregations: []Aggregation{AggRawJson(`{ "type" : "count", "name" : "count" }`), AggLongSum("revenue", "dimension_sum"), AggHyperUnique("unique_devices", "unique_devices")},
			PostAggregations: []PostAggregation{
				PostAggArithmetic("Revenue/Event", "/",
					[]PostAggregation{
						PostAggFieldAccessor("revenue"),
						PostAggRawJson(`{ "type" : "fieldAccess", "fieldName" : "count" }`)})},
		}

		client := Client{
			Url:   broker,
			Debug: true,
		}

		err := client.Query(query)
		fmt.Println("request", client.LastRequest)
		So(err, ShouldEqual, nil)

		fmt.Println("response", client.LastResponse)

		fmt.Printf("query.QueryResult:\n%v", query.QueryResult)

	})
}

func TestSearch(t *testing.T) {
	Convey("TestSearch", t, func() {
		broker, found := syscall.Getenv("DruidBroker")
		appId, found2 := syscall.Getenv("TestAppId")
		So(found && found2, ShouldEqual, true)

		query := &QuerySearch{
			DataSource:       "events_agg",
			Intervals:        []string{"2016-05-01T00:00/2016-05-01T01"},
			Granularity:      GranAll,
			Filter:           FilterSelector("app_id", appId),
			SearchDimensions: []string{"attribution_campaign_id", "hour"},
			Query:            SearchQueryInsensitiveContains("131"),
			Sort:             SearchSortLexicographic,
		}

		client := Client{
			Url:   broker,
			Debug: true,
		}

		err := client.Query(query)
		So(err, ShouldEqual, nil)

		fmt.Println("request", client.LastRequest)
		fmt.Println("response", client.LastResponse)

		fmt.Printf("query.QueryResult:\n%v", query.QueryResult)

	})
}

func TestTopN(t *testing.T) {
	Convey("TestTopN", t, func() {
		broker, found := syscall.Getenv("DruidBroker")
		appId, found2 := syscall.Getenv("TestAppId")
		So(found && found2, ShouldEqual, true)

		query := &QueryTopN{
			DataSource:   "events_agg",
			Intervals:    []string{"2016-05-01T00:00/2016-05-01T01"},
			Granularity:  GranAll,
			Filter:       FilterAnd(FilterSelector("app_id", appId), FilterSelector("attribution_network_key", ""), nil),
			Dimension:    "device_os",
			Threshold:    50,
			Metric:       TopNMetricNumeric("Revenue/Event"),
			Aggregations: []Aggregation{AggRawJson(`{ "type" : "count", "name" : "count" }`), AggLongSum("revenue", "dimension_sum"), AggHyperUnique("unique_devices", "unique_devices")},
			PostAggregations: []PostAggregation{
				PostAggArithmetic("Revenue/Event", "/",
					[]PostAggregation{
						PostAggFieldAccessor("revenue"),
						PostAggRawJson(`{ "type" : "fieldAccess", "fieldName" : "count" }`)})},
		}

		client := Client{
			Url:   broker,
			Debug: true,
		}

		err := client.Query(query)
		fmt.Println("request", client.LastRequest)
		So(err, ShouldEqual, nil)

		fmt.Println("response", client.LastResponse)

		fmt.Printf("query.QueryResult:\n%v", query.QueryResult)

	})
}

func TestTimeseries(t *testing.T) {
	Convey("TestTimeseries", t, func() {
		broker, found := syscall.Getenv("DruidBroker")
		appId, found2 := syscall.Getenv("TestAppId")
		So(found && found2, ShouldEqual, true)

		query := &QueryTimeseries{
			DataSource:   "events_agg",
			Intervals:    []string{"2016-05-01T00:00/2016-05-01T05"},
			Granularity:  GranPeriod("PT1H"),
			Filter:       FilterAnd(FilterSelector("app_id", appId), FilterSelector("attribution_network_key", ""), nil),
			Aggregations: []Aggregation{AggRawJson(`{ "type" : "count", "name" : "count" }`), AggLongSum("revenue", "dimension_sum"), AggHyperUnique("unique_devices", "unique_devices")},
			PostAggregations: []PostAggregation{
				PostAggArithmetic("Revenue/Event", "/", []PostAggregation{
					PostAggFieldAccessor("revenue"),
					PostAggRawJson(`{ "type" : "fieldAccess", "fieldName" : "count" }`)}),
				PostAggArithmetic("Revenue/User", "/", []PostAggregation{
					PostAggFieldAccessor("revenue"),
					PostAggFieldHyperUnique("unique_devices")})},
		}

		client := Client{
			Url:   broker,
			Debug: true,
		}

		err := client.Query(query)
		fmt.Println("request", client.LastRequest)
		So(err, ShouldEqual, nil)

		fmt.Println("response", client.LastResponse)

		fmt.Printf("query.QueryResult:\n%v", query.QueryResult)

	})
}

func TestSelect(t *testing.T) {
	Convey("TestSelect", t, func() {
		broker, found := syscall.Getenv("DruidBroker")
		appId, found2 := syscall.Getenv("TestAppId")
		So(found && found2, ShouldEqual, true)

		query := &QuerySelect{
			DataSource:  "events",
			Intervals:   []string{"2016-05-01T00:00/2016-05-01T01:00"},
			Granularity: GranAll,
			Filter:      FilterAnd(FilterSelector("app_id", appId), FilterSelector("attribution_network_key", ""), nil),
			PagingSpec:  PagingSpec{PagingIdentifiers: PagingIdEmpty{}, Threshold: 5},
		}

		client := Client{
			Url:   broker,
			Debug: true,
		}

		err := client.Query(query)
		fmt.Println("request", client.LastRequest)
		So(err, ShouldEqual, nil)

		fmt.Println("response", client.LastResponse)

		fmt.Printf("query.QueryResult:\n%v", query.QueryResult)

	})
}

func TestTheta(t *testing.T) {
	Convey("TestTheta", t, func() {
		broker, found := syscall.Getenv("DruidBroker")
		appId, found2 := syscall.Getenv("TestAppId")
		So(found && found2, ShouldEqual, true)

		query := &QueryGroupBy{
			DataSource:       "events_agg",
			Intervals:        []string{"2016-05-01T00:00/2016-05-02T00:00"},
			Granularity:      GranAll,
			Filter:           FilterAnd(FilterSelector("app_id", appId), FilterOr(FilterSelector("event_name", "Purchase"), FilterSelector("event_name", "_Install"))),
			Aggregations:     []Aggregation{AggFiltered(*FilterSelector("event_name", "_Install"), AggThetaSketch("Installs", "theta_devices")), AggFiltered(*FilterSelector("event_name", "Purchase"), AggThetaSketch("Purchases", "theta_devices"))},
			PostAggregations: []PostAggregation{PostAggThetaOp("InstallAndPurchase", ThetaIntersect, []PostAggregation{PostAggFieldAccessor("Installs"), PostAggFieldAccessor("Purchases")}), PostAggThetaEstimate("InstallAndPurchaseEstimate", PostAggFieldAccessor("InstallAndPurchase"))},
		}

		client := Client{
			Url:   broker,
			Debug: true,
		}

		err := client.Query(query)
		fmt.Println("request", client.LastRequest)
		So(err, ShouldEqual, nil)

		fmt.Println("response", client.LastResponse)

		fmt.Printf("query.QueryResult:\n%v", query.QueryResult)

	})
}

func TestHaving(t *testing.T) {
	Convey("TestHaving", t, func() {
		broker, found := syscall.Getenv("DruidBroker")
		appId, found2 := syscall.Getenv("TestAppId")
		So(found && found2, ShouldEqual, true)

		query := &QueryGroupBy{
			DataSource:   "events_agg",
			Dimensions:   []DimSpec{"attribution_network"},
			Intervals:    []string{"2016-05-01T00:00/2016-05-05T00:00"},
			Granularity:  GranDuration(3600000),
			Filter:       FilterAnd(FilterSelector("app_id", appId), FilterRegex("attribution_network", "^[SsBbRr]")),
			Aggregations: []Aggregation{AggDoubleMax("maximum", "dimension_sum"), AggDoubleSum("total", "dimension_sum"), AggDoubleSum("count", "count")},
			Having:       HavingAnd(HavingGreaterThan("maximum", 10000), HavingLessThan("total", 500000)),
		}

		client := Client{
			Url:   broker,
			Debug: true,
		}

		err := client.Query(query)
		fmt.Println("request", client.LastRequest)
		So(err, ShouldEqual, nil)

		fmt.Println("response", client.LastResponse)

		fmt.Printf("query.QueryResult:\n%v", query.QueryResult)

	})
}

func TestLookup(t *testing.T) {
	Convey("TestLookup", t, func() {
		broker, found := syscall.Getenv("DruidBroker")
		appId, found2 := syscall.Getenv("TestAppId")
		So(found && found2, ShouldEqual, true)

		query := &QueryTopN{
			DataSource:   "events_agg",
			Intervals:    []string{"2016-05-01T00:00/2016-05-03T01"},
			Granularity:  GranHour,
			Filter:       FilterAnd(FilterSelector("app_id", appId), FilterIn("attribution_days_since", "0", "1"), FilterSelector("event_name", "Purchase")),
			Dimension:    DimExtraction("attribution_hours_since", "since", DimExFnLookup(LookupMap(map[string]string{"1": "1", "2": "1", "3": "2", "4": "2", "5": "10"}, IsOneToOne(false)), RetainMissingValue(false), ReplaceMissingValueWith("Somethin' else"))),
			Threshold:    50,
			Metric:       TopNMetricAlphaNumeric(),
			Aggregations: []Aggregation{AggCount("count"), AggLongSum("revenue", "dimension_sum"), AggHyperUnique("unique_devices", "unique_devices")},
			PostAggregations: []PostAggregation{
				PostAggArithmetic("Revenue/Event", "/",
					[]PostAggregation{
						PostAggFieldAccessor("revenue"),
						PostAggRawJson(`{ "type" : "fieldAccess", "fieldName" : "count" }`)})},
		}

		client := Client{
			Url:   broker,
			Debug: true,
		}

		err := client.Query(query)
		fmt.Println("request", client.LastRequest)
		So(err, ShouldEqual, nil)

		fmt.Println("response", client.LastResponse)

		fmt.Printf("query.QueryResult:\n%v", query.QueryResult)

	})
}
