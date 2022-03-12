package plugin

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func TestEpochGroupSecondsShouldBeReplacedInTheFinalQueryForTables(t *testing.T) {
	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value INTEGER);
		INSERT INTO test(time, value)
		VALUES (4, 1), (13, 2), (34, 4);
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText:   `SELECT $__unixEpochGroupSeconds("time", 10) as window, value FROM test`,
		TimeColumns: []string{},
	})
	dataQuery.QueryType = tableType

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error != nil {
		t.Errorf("Unexpected error - %s", response.Error)
	}

	if len(response.Frames) != 1 {
		t.Errorf(
			"Expected one frame but got - %d: Frames %+v", len(response.Frames), response.Frames,
		)
	}

	expectedFrame := data.NewFrame(
		"response",
		data.NewField("window", nil, []*int64{intPointer(0), intPointer(10), intPointer(30)}),
		data.NewField("value", nil, []*int64{intPointer(1), intPointer(2), intPointer(4)}),
	)
	expectedFrame.Meta = &data.FrameMeta{
		ExecutedQueryString: `SELECT cast(("time" / 10) as int) * 10 as window, value FROM test`,
	}

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Error(diff)
	}
}

func TestEpochGroupSecondsShouldFillInNullValuesForTables(t *testing.T) {
	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value INTEGER);
		INSERT INTO test(time, value)
		VALUES (4, 1), (13, 2), (34, 4);
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText:   `SELECT $__unixEpochGroupSeconds("time", 10, NULL) as window, value FROM test`,
		TimeColumns: []string{"window"},
	})
	dataQuery.QueryType = tableType

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error != nil {
		t.Errorf("Unexpected error - %s", response.Error)
	}

	if len(response.Frames) != 1 {
		t.Fatalf(
			"Expected one frame but got - %d: Frames %+v", len(response.Frames), response.Frames,
		)
	}

	expectedFrame := data.NewFrame(
		"response",
		data.NewField(
			"window", nil, []*time.Time{
				unixTimePointer(0),
				unixTimePointer(10),
				unixTimePointer(20),
				unixTimePointer(30),
			},
		),
		data.NewField("value", nil, []*int64{intPointer(1), intPointer(2), nil, intPointer(4)}),
	)
	expectedFrame.Meta = &data.FrameMeta{
		ExecutedQueryString: `SELECT cast(("time" / 10) as int) * 10 as window, value FROM test`,
	}

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Error(diff)
	}
}

func TestEpochGroupSecondsShouldFillInNullValuesForTimeSeriesWithDoubleGaps(t *testing.T) {
	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value INTEGER);
		INSERT INTO test(time, value)
		VALUES (4, 1), (13, 2), (44, 5);
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText:   `SELECT $__unixEpochGroupSeconds("time", 10, NULL) as window, value FROM test`,
		TimeColumns: []string{"window"},
	})
	dataQuery.QueryType = timeSeriesType

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error != nil {
		t.Errorf("Unexpected error - %s", response.Error)
	}

	if len(response.Frames) != 1 {
		t.Fatalf(
			"Expected one frame but got - %d: Frames %+v", len(response.Frames), response.Frames,
		)
	}

	expectedFrame := data.NewFrame(
		"value",
		data.NewField(
			"window",
			nil,
			[]*time.Time{
				unixTimePointer(0),
				unixTimePointer(10),
				unixTimePointer(20),
				unixTimePointer(30),
				unixTimePointer(40),
			},
		),
		data.NewField(
			"value", nil, []*int64{intPointer(1), intPointer(2), nil, nil, intPointer(5)},
		),
	)
	expectedFrame.Meta = &data.FrameMeta{
		ExecutedQueryString: `SELECT cast(("time" / 10) as int) * 10 as window, value FROM test`,
	}

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Error(diff)
	}
}

func TestEpochGroupSecondsShouldFillInNullValuesWithMultipleTimeColumns(t *testing.T) {
	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, other_time INTEGER, value INTEGER);
		INSERT INTO test(time, other_time, value)
		VALUES (4, 11, 1), (13, 12, 2), (34, 13, 4);
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText: `
			SELECT
				$__unixEpochGroupSeconds("time", 10, NULL) as window
				, other_time
				, value
			 FROM test
		`,
		TimeColumns: []string{"window", "other_time"},
	})
	dataQuery.QueryType = tableType

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error != nil {
		t.Errorf("Unexpected error - %s", response.Error)
	}

	if len(response.Frames) != 1 {
		t.Errorf(
			"Expected one frame but got - %d: Frames %+v", len(response.Frames), response.Frames,
		)
	}

	expectedFrame := data.NewFrame(
		"response",
		data.NewField("window", nil, []*time.Time{
			unixTimePointer(0),
			unixTimePointer(10),
			unixTimePointer(20),
			unixTimePointer(30),
		}),
		data.NewField("other_time", nil, []*time.Time{
			unixTimePointer(11),
			unixTimePointer(12),
			nil,
			unixTimePointer(13),
		}),
		data.NewField("value", nil, []*int64{intPointer(1), intPointer(2), nil, intPointer(4)}),
	)
	expectedFrame.Meta = &data.FrameMeta{
		// we test this content elsewhere and do not care about it
		ExecutedQueryString: response.Frames[0].Meta.ExecutedQueryString,
	}

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Error(diff)
	}
}

func TestEpochGroupSecondsShouldBeReplacedInTheFinalQueryForTimeSeries(t *testing.T) {
	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value INTEGER);
		INSERT INTO test(time, value)
		VALUES (4, 1), (13, 2), (34, 4);
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText:   `SELECT $__unixEpochGroupSeconds("time", 10) as window, value FROM test`,
		TimeColumns: []string{"window"},
	})
	dataQuery.QueryType = timeSeriesType

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error != nil {
		t.Errorf("Unexpected error - %s", response.Error)
	}

	if len(response.Frames) != 1 {
		t.Errorf(
			"Expected one frame but got - %d: Frames %+v", len(response.Frames), response.Frames,
		)
	}

	expectedFrame := data.NewFrame(
		"value",
		data.NewField("window", nil, []*time.Time{
			timePointer(time.Unix(0, 0)),
			timePointer(time.Unix(10, 0)),
			timePointer(time.Unix(30, 0)),
		}),
		data.NewField("value", nil, []*int64{intPointer(1), intPointer(2), intPointer(4)}),
	)
	expectedFrame.Meta = &data.FrameMeta{
		ExecutedQueryString: `SELECT cast(("time" / 10) as int) * 10 as window, value FROM test`,
	}

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Error(diff)
	}
}

func TestEpochGroupSecondsWithMultiFrameTimeseriesAndGaps(t *testing.T) {
	var inputFrame *data.Frame
	mockableLongToWide = func(a *data.Frame, b *data.FillMissing) (*data.Frame, error) {
		inputFrame = a
		return data.LongToWide(a, b)
	}

	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value REAL, name TEXT);
		INSERT INTO test(time, value, name)
		VALUES (11, 11.1, 'one'), (32, 22.2, 'two');
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText: `
			SELECT $__unixEpochGroupSeconds("time", 10, NULL) as window, name, value
			FROM test GROUP BY 1, name ORDER BY 1 ASC
		`,
		TimeColumns: []string{"window"},
	})
	dataQuery.QueryType = "time series"

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error != nil {
		t.Errorf("Unexpected error - %s", response.Error)
	}

	expectedInputFrame := data.NewFrame(
		"response",
		data.NewField("window", nil, []*time.Time{
			unixTimePointer(10), unixTimePointer(20), unixTimePointer(30),
		}),
		data.NewField("name", nil, []*string{strPointer("one"), nil, strPointer("two")}),
		data.NewField("value", nil, []*float64{floatPointer(11.1), nil, floatPointer(22.2)}),
	)
	// we use the response as we do not care about the value (tested elsewhere)
	expectedInputFrame.Meta = &data.FrameMeta{
		ExecutedQueryString: response.Frames[0].Meta.ExecutedQueryString, Type: data.FrameTypeTimeSeriesWide,
	}

	if diff := cmp.Diff(expectedInputFrame, inputFrame, cmpOption...); diff != "" {
		t.Error("Unexpected input frame into the time series conversion")
		t.Fatal(diff)
	}

	if len(response.Frames) != 2 {
		t.Errorf(
			"Expected two frames but got - %d: Frames %+v", len(response.Frames), response.Frames,
		)
	}

	expectedOutputFrames := make([]*data.Frame, 2)
	expectedOutputFrames[0] = data.NewFrame(
		"value one",
		data.NewField("window", nil, []time.Time{
			time.Unix(10, 0), time.Unix(20, 0), time.Unix(30, 0)},
		),
		data.NewField(
			"value",
			map[string]string{"name": "one"},
			[]*float64{floatPointer(11.1), nil, nil},
		),
	)
	// we use the response as we do not care about the value (tested elsewhere)
	expectedOutputFrames[0].Meta = response.Frames[0].Meta

	expectedOutputFrames[1] = data.NewFrame(
		"value two",
		data.NewField("window", nil, []time.Time{
			time.Unix(10, 0), time.Unix(20, 0), time.Unix(30, 0)},
		),
		data.NewField(
			"value",
			map[string]string{"name": "two"},
			[]*float64{nil, nil, floatPointer(22.2)},
		),
	)
	// we use the response as we do not care about the value (tested elsewhere)
	expectedOutputFrames[1].Meta = response.Frames[1].Meta

	for idx, frame := range response.Frames {
		if diff := cmp.Diff(expectedOutputFrames[idx], frame, cmpOption...); diff != "" {
			t.Error("Unexpected output frames")
			t.Error(diff)
		}
	}
}

func TestEpochGroupSecondsShouldNotAcceptOneArgument(t *testing.T) {
	dataQuery := getDataQuery(queryModel{
		QueryText:   `SELECT $__unixEpochGroupSeconds("time") as window, value FROM test`,
		TimeColumns: []string{},
	})
	dataQuery.QueryType = tableType

	response := query(dataQuery, pluginConfig{Path: "dbPath"})
	if response.Error == nil {
		t.Errorf("Expected error but got nothing. Response: %+v", response)
	}

	if !strings.Contains(fmt.Sprintf("%v", response.Error), "unsupported number of arguments") {
		t.Errorf("Expected argument error but got: %+v", response.Error)
	}
}

func TestEpochGroupSecondsShouldNotAcceptAnyStringAsGapValue(t *testing.T) {
	dataQuery := getDataQuery(queryModel{
		QueryText:   `SELECT $__unixEpochGroupSeconds("time", 10, something) as window, value FROM test`,
		TimeColumns: []string{"window"},
	})
	dataQuery.QueryType = tableType

	response := query(dataQuery, pluginConfig{Path: "dbPath"})
	if response.Error == nil {
		t.Errorf("Expected error but got nothing. Response: %+v", response)
	}

	if !strings.Contains(fmt.Sprintf("%v", response.Error), "unsupported gap filling value") {
		t.Errorf("Expected argument error but got: %+v", response.Error)
	}
}

func TestEpochGroupSecondsShouldRequireATimeColumnForGapFilling(t *testing.T) {
	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value INTEGER);
		INSERT INTO test(time, value)
		VALUES (4, 1), (13, 2), (34, 4);
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText:   `SELECT $__unixEpochGroupSeconds("time", 10, NULL) as window, value FROM test`,
		TimeColumns: []string{"not_found_column"},
	})
	dataQuery.QueryType = tableType

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error == nil {
		t.Errorf("Expected error but got nothing. Response: %+v", response)
	}

	if !strings.Contains(fmt.Sprintf("%v", response.Error), "no time column found") {
		t.Errorf("Expected no time column error but got: %+v", response.Error)
	}
}

func TestEpochGroupSecondsShouldRequireAnOrderedTimeColumnForGapFilling(t *testing.T) {
	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value INTEGER);
		INSERT INTO test(time, value)
		VALUES (14, 1), (23, 2), (4, 4);
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText:   `SELECT $__unixEpochGroupSeconds("time", 10, NULL) as window, value FROM test`,
		TimeColumns: []string{"window"},
	})
	dataQuery.QueryType = tableType

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error == nil {
		t.Errorf("Expected error but got nothing. Response: %+v", response)
	}

	if !strings.Contains(
		fmt.Sprintf("%v", response.Error), "unordered time value",
	) {
		t.Errorf("Expected unordered time column error but got: %+v", response.Error)
	}
}

func TestEpochGroupSecondsShouldRequireNonNullTimeColumnForGapFilling(t *testing.T) {
	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value INTEGER);
		INSERT INTO test(time, value)
		VALUES (14, 1), (NULL, 2), (24, 4);
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText:   `SELECT $__unixEpochGroupSeconds("time", 10, NULL) as window, value FROM test`,
		TimeColumns: []string{"window"},
	})
	dataQuery.QueryType = tableType

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error == nil {
		t.Errorf("Expected error but got nothing. Response: %+v", response)
	}

	if !strings.Contains(
		fmt.Sprintf("%v", response.Error), "NULL value in time column",
	) {
		t.Errorf("Expected null in time column error but got: %+v", response.Error)
	}
}
