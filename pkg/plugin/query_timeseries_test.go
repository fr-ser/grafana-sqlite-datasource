package plugin

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func TestIgnoreNonTimeSeriesQuery(t *testing.T) {
	var longToWideCalled bool
	mockableLongToWide = func(a *data.Frame, b *data.FillMissing) (*data.Frame, error) {
		longToWideCalled = true
		return data.NewFrame("response"), nil
	}

	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value REAL, name TEXT);
		INSERT INTO test(time, value, name)
		VALUES (21, 21.1, 'one'), (22.0, 22.2, 'two'), (23, 23.3, 'three');
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText: "SELECT * FROM test", TimeColumns: []string{"time"},
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

	if longToWideCalled {
		t.Errorf("Expected to not call 'longToWide' for non time series queries")
	}

	expectedFrame := data.NewFrame(
		"response",
		data.NewField("time", nil, []*time.Time{
			timePointer(time.Unix(21, 0)),
			timePointer(time.Unix(22, 0)),
			timePointer(time.Unix(23, 0)),
		}),
		data.NewField("value", nil, []*float64{
			floatPointer(21.1), floatPointer(22.2), floatPointer(23.3),
		}),
		data.NewField("name", nil, []*string{
			strPointer("one"), strPointer("two"), strPointer("three"),
		}),
	)
	expectedFrame.Meta = &data.FrameMeta{ExecutedQueryString: "SELECT * FROM test"}

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Error(diff)
	}
}

func TestIgnoreWideTimeSeriesQuery(t *testing.T) {
	var longToWideCalled bool
	mockableLongToWide = func(a *data.Frame, b *data.FillMissing) (*data.Frame, error) {
		longToWideCalled = true
		return data.NewFrame("response"), nil
	}

	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value REAL);
		INSERT INTO test(time, value)
		VALUES (21, 21.1), (22.0, 22.2), (23, 23.3);
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText: "SELECT * FROM test", TimeColumns: []string{"time"},
	})
	dataQuery.QueryType = "time series"

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error != nil {
		t.Errorf("Unexpected error - %s", response.Error)
	}

	if len(response.Frames) != 1 {
		t.Errorf(
			"Expected one frame but got - %d: Frames %+v", len(response.Frames), response.Frames,
		)
	}

	if longToWideCalled {
		t.Errorf("Expected to not call 'longToWide' for wide time series queries")
	}

	expectedFrame := data.NewFrame(
		"value",
		data.NewField("time", nil, []*time.Time{
			timePointer(time.Unix(21, 0)),
			timePointer(time.Unix(22, 0)),
			timePointer(time.Unix(23, 0)),
		}),
		data.NewField("value", nil, []*float64{
			floatPointer(21.1), floatPointer(22.2), floatPointer(23.3),
		}),
	)
	expectedFrame.Meta = &data.FrameMeta{ExecutedQueryString: "SELECT * FROM test"}

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Error(diff)
	}
}

func TestConvertLongTimeSeriesQuery(t *testing.T) {
	var inputFrame *data.Frame
	mockableLongToWide = func(a *data.Frame, b *data.FillMissing) (*data.Frame, error) {
		inputFrame = a
		return data.LongToWide(a, b)
	}

	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value REAL, name TEXT);
		INSERT INTO test(time, value, name)
		VALUES (21, 21.1, 'one'), (22, 22.2, 'two');
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText: "SELECT * FROM test", TimeColumns: []string{"time"},
	})
	dataQuery.QueryType = "time series"

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error != nil {
		t.Errorf("Unexpected error - %s", response.Error)
	}

	expectedInputFrame := data.NewFrame(
		"response",
		data.NewField("time", nil, []*time.Time{
			timePointer(time.Unix(21, 0)), timePointer(time.Unix(22, 0)),
		}),
		data.NewField("value", nil, []*float64{floatPointer(21.1), floatPointer(22.2)}),
		data.NewField("name", nil, []*string{strPointer("one"), strPointer("two")}),
	)
	expectedInputFrame.Meta = &data.FrameMeta{
		Type: data.FrameTypeTimeSeriesWide, ExecutedQueryString: "SELECT * FROM test",
	}

	if diff := cmp.Diff(expectedInputFrame, inputFrame, cmpOption...); diff != "" {
		t.Error("Unexpected input frame into the time series conversion")
		t.Error(diff)
	}

	if len(response.Frames) != 2 {
		t.Errorf(
			"Expected two frames but got - %d: Frames %+v", len(response.Frames), response.Frames,
		)
	}

	expectedOutputFrames := make([]*data.Frame, 2)
	expectedOutputFrames[0] = data.NewFrame(
		"value one",
		data.NewField("time", nil, []time.Time{time.Unix(21, 0), time.Unix(22, 0)}),
		data.NewField(
			"value",
			map[string]string{"name": "one"},
			[]*float64{floatPointer(21.1), nil},
		),
	)
	expectedOutputFrames[0].Meta = &data.FrameMeta{ExecutedQueryString: "SELECT * FROM test"}

	expectedOutputFrames[1] = data.NewFrame(
		"value two",
		data.NewField("time", nil, []time.Time{time.Unix(21, 0), time.Unix(22, 0)}),
		data.NewField(
			"value",
			map[string]string{"name": "two"},
			[]*float64{nil, floatPointer(22.2)},
		),
	)
	expectedOutputFrames[1].Meta = &data.FrameMeta{ExecutedQueryString: "SELECT * FROM test"}

	for idx, frame := range response.Frames {
		if diff := cmp.Diff(expectedOutputFrames[idx], frame, cmpOption...); diff != "" {
			t.Error("Unexpected output frames")
			t.Error(diff)
		}
	}
}
