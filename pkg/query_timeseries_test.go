package main

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func TestIgnoreNonTimeseriesQuery(t *testing.T) {
	var longToWideCalled bool
	mockableLongToWide = func(a *data.Frame, b *data.FillMissing) (*data.Frame, error) {
		longToWideCalled = true
		return a, nil
	}

	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value REAL, name TEXT);
		INSERT INTO test(time, value, name)
		VALUES (21, 21.1, 'one'), (22.3, 22.2, 'two'), (23, 23.3, 'three');
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText: "SELECT * FROM test", TimeColumns: []string{"time"},
	})
	dataQuery.QueryType = "table"

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
		t.Errorf("Expected to not call 'longToWide' for non timeseries queries")
	}

	expectedFrame := data.NewFrame(
		"response",
		data.NewField("time", nil, []*time.Time{
			timePointer(time.Unix(21, 0)), timePointer(time.Unix(22, 0)),
			timePointer(time.Unix(23, 0)),
		}),
		data.NewField("value", nil, []*float64{
			floatPointer(21.1), floatPointer(22.2), floatPointer(23.3),
		}),
		data.NewField("name", nil, []*string{
			strPointer("one"), strPointer("two"), strPointer("three"),
		}),
	)

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Error(diff)
	}
}

func TestConvertTimeseriesQuery(t *testing.T) {
	var inputFrame *data.Frame
	outputFrame := data.NewFrame("response")
	mockableLongToWide = func(a *data.Frame, b *data.FillMissing) (*data.Frame, error) {
		inputFrame = a
		return outputFrame, nil
	}

	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value REAL, name TEXT);
		INSERT INTO test(time, value, name)
		VALUES (21, 21.1, 'one'), (21.3, 22.2, 'two'), (21, 23.3, 'three');
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText: "SELECT * FROM test", TimeColumns: []string{"time"},
	})
	dataQuery.QueryType = "timeseries"

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error != nil {
		t.Errorf("Unexpected error - %s", response.Error)
	}

	expectedInputFrame := data.NewFrame(
		"response",
		data.NewField("time", nil, []*time.Time{
			timePointer(time.Unix(21, 0)), timePointer(time.Unix(21, 0)),
			timePointer(time.Unix(21, 0)),
		}),
		data.NewField("value", nil, []*float64{
			floatPointer(21.1), floatPointer(22.2), floatPointer(23.3),
		}),
		data.NewField("name", nil, []*string{
			strPointer("one"), strPointer("two"), strPointer("three"),
		}),
	)

	if diff := cmp.Diff(expectedInputFrame, inputFrame, cmpOption...); diff != "" {
		t.Error(diff)
	}

	if len(response.Frames) != 1 {
		t.Errorf(
			"Expected one frame but got - %d: Frames %+v", len(response.Frames), response.Frames,
		)
	}

	if diff := cmp.Diff(outputFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Error("Did not receive expected outputFrame from longToWide call")
		t.Error(diff)

	}
}
